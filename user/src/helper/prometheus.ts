import prometheus from "prom-client";

export const register = new prometheus.Registry();

export const histogram: prometheus.Histogram = new prometheus.Histogram({
    name: 'user_histogram',
    help: 'Histogram about User Request',
    labelNames: ['method', 'code', 'type'],
    buckets: [0.1, 0.15, 0.2, 0.25, 0.3 ]
})

export const counter: prometheus.Counter = new prometheus.Counter({
    name: 'user_counter',
    help: 'Counter about user request',
    labelNames: ['type']
})

export const totalCPU: prometheus.Gauge = new prometheus.Gauge({
    name: 'total_cpu',
    help: 'Total CPU on user services',
    labelNames: ['version']
})

export const totalMemory: prometheus.Gauge = new prometheus.Gauge({
    name: 'total_memory',
    help: 'Total Memory on user services',
    labelNames: ['version']
})