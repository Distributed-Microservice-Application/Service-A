import http from 'k6/http';
import { check, sleep } from 'k6';

// Options for the test
export const options = {
    scenarios: {
        basic: {
            executor: 'constant-vus',
            vus: 5,
            duration: '10s',
        },
        load: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '1m', target: 15 },
                { duration: '1m', target: 15 },
                { duration: '30s', target: 0 },
            ],
            gracefulRampDown: '5s',
        },
    },
    thresholds: {
        'checks': ['rate>0.95'],
        'http_req_duration': ['p(95)<600'],
    },
};

const Rand = (min, max) => {
    return Math.floor(Math.random() * (max - min + 1)) + min;
};

export default function () {
    const A = Rand(1, 500);
    const B = Rand(1, 300);

    const payload = JSON.stringify({
        a: A,
        b: B,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Call the Nginx load balancer endpoint
    const res = http.post('http://localhost:8090/sum', payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'result is correct': (r) => {
            const body = JSON.parse(r.body);
            return body.result === (A + B);
        },
    });

    sleep(1);
}
