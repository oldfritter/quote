nohup ./cmd/api >> logs/api.log &
nohup ./schedules/schedule >> logs/schedule.log &
nohup ./workers/worker >> logs/worker.log &
nohup ./cmd/getData >> logs/getData.log &
