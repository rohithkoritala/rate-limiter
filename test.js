import http from 'k6/http';
import { check } from 'k6';
import { sleep } from 'k6';

export let options = {
  vus: 500, // virtual users
  duration: '30s',
};

export default function () {
  const key = `user${Math.floor(Math.random() * 10000)}`; // Random key per user
  const res = http.get(`http://localhost:8080/check?key=${key}`);
  
  check(res, {
    'status is 200 or 429': (r) => r.status === 200 || r.status === 429,
    'response body is valid': (r) => r.body === 'allowed' || r.body === 'rate limit exceeded',
  });

  sleep(0.01); // slight delay to avoid overwhelming the machine
}
