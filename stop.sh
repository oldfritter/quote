cat pids/api.pid  | xargs kill -INT
cat pids/getData.pid  | xargs kill -INT
cat pids/schedule.pid  | xargs kill -INT
cat pids/worker.pid  | xargs kill -INT
