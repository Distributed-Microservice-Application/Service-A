import grpc from 'k6/net/grpc';
import { check } from 'k6';

const client = new grpc.Client();
client.load(['../proto'], 'summation.proto');

// Options for the test
export const options = {
    scenarios: {
        basic: {
            executor: 'constant-vus',
            vus: 5,
            duration: '10s',
        },
        load: {
            /*
                Ramping up the number of virtual users (VUs) over time.
                - Start with 0 VUs.
                - Ramp up to 15 VUs over 1 minute.
                - Maintain 15 VUs for another minute.
                - Ramp down to 0 VUs over 30 seconds.
                - Gracefully ramp down over 5 seconds.
            */
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '1m', target: 15 },
                { duration: '1m', target: 15 },
                { duration: '30s', target: 0 },
            ],
            // This is the time to wait before starting the ramp down phase.
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
    // Connect to the gRPC server
    client.connect('localhost:50051', { plaintext: true });

    const A = Rand(1, 500);
    const B = Rand(1, 300);

    const request = {
        a: A,
        b: B,
    };

    const res = A + B;
    // Use the correct service and method name from the proto file
    const response = client.invoke('summation.SummationService/CalculateSum', request);

   check(response, {
      'status is OK': (r) => r && r.status === grpc.StatusOK,
      'result is correct': (r) => {
           return r && r.message && r.message.result === res;
      },
   });

   // close the connection after the request
   client.close();
}