import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

const realErrors = new Counter('real_errors');

export const options = {
  vus: 10,
  duration: '30s',
  thresholds: {
    http_req_duration: ['p(95) < 300'],
    real_errors: ['count < 3'],
  },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  // Уникальные ID на основе времени и VU для избежания коллизий
  const uniqueId = Date.now().toString() + '_' + __VU;
  const teamName = 'team_' + uniqueId;
  const authorId = 'author_' + uniqueId;
  const reviewer1Id = 'rev1_' + uniqueId;
  const reviewer2Id = 'rev2_' + uniqueId;
  const prId = 'pr_' + uniqueId;

  // Создание команды (3 участника: автор + 2 ревьювера)
  const teamRes = http.post(`${BASE_URL}/team/add`, JSON.stringify({
    team_name: teamName,
    members: [
      { user_id: authorId, username: "Author", is_active: true },
      { user_id: reviewer1Id, username: "Reviewer1", is_active: true },
      { user_id: reviewer2Id, username: "Reviewer2", is_active: true }
    ]
  }), { headers: { 'Content-Type': 'application/json' } });

  // Проверка создания команды, при ошибке думаем, что проблема реальная
  if (!check(teamRes, { 'team created': (r) => r.status === 201 })) {
    realErrors.add(1);
    return;
  }

  // Создание PR
  const prRes = http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify({
    pull_request_id: prId,
    pull_request_name: "Load test PR",
    author_id: authorId
  }), { headers: { 'Content-Type': 'application/json' } });

  if (!check(prRes, { 'PR created': (r) => r.status === 201 })) {
    realErrors.add(1);
    return;
  }

  let prData;
  try {
    prData = JSON.parse(prRes.body);
  } catch (e) {
    realErrors.add(1);
    return;
  }

  if (!prData.pr) {
    realErrors.add(1);
    return;
  }

  const assignedReviewers = prData.pr.assigned_reviewers || [];

  if (assignedReviewers.length > 0) {
    const reassignRes = http.post(`${BASE_URL}/pullRequest/reassign`, JSON.stringify({
      pull_request_id: prId,
      old_reviewer_id: assignedReviewers[0]
    }), { headers: { 'Content-Type': 'application/json' } });

    const isReassignValid = reassignRes.status === 200 || reassignRes.status === 409;
    if (!isReassignValid) {
      realErrors.add(1);
    }
  }

  // Получение списка PR у ревьювера
  const reviewsRes = http.get(`${BASE_URL}/users/getReview?user_id=${reviewer1Id}`);
  if (!check(reviewsRes, { 'get reviews': (r) => r.status === 200 })) {
    realErrors.add(1);
  }

  // Мерж PR
  const mergeRes = http.post(`${BASE_URL}/pullRequest/merge`, JSON.stringify({
    pull_request_id: prId
  }), { headers: { 'Content-Type': 'application/json' } });

  if (!check(mergeRes, { 'PR merged': (r) => r.status === 200 })) {
    realErrors.add(1);
  }

  // Запрос статистики
  const statsRes = http.get(`${BASE_URL}/stats/reviewer-assignments`);
  if (!check(statsRes, { 'stats loaded': (r) => r.status === 200 })) {
    realErrors.add(1);
  }

  // 1 запрос/сек на VU → 10 RPS суммарно
  sleep(1);
}