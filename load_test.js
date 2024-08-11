import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 2500 }, // Ramp-up to 100 users over 30 seconds
        { duration: '1m', target: 5000 },  // Hold at 200 users for 1 minute
        { duration: '30s', target: 0 },   // Ramp-down to 0 users over 30 seconds
    ],
};

export default function () {
    let url = 'http://kingoftheheap.dev/shorten';
    let payload = JSON.stringify({
        url: 'https://www.example.com/' + Math.random().toString(36).substring(7)
    });

    let params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let res = http.post(url, payload, params);
    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time is < 400ms': (r) => r.timings.duration < 400,
    });

    sleep(1); // Simulate some delay between requests
}
