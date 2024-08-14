import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 10000}, // ramp up
        { duration: '1m', target: 20000}, //hold
        { duration: '30s', target: 0}, //ramp down
        
    ],
};

export default function () {
    let url = 'https://skny.link/shorten';
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

    sleep(1); // delay
}
