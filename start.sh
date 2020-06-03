nohup ./cmd/api >> logs/api.log &
nohup ./cmd/schedule >> logs/schedule.log &
nohup ./cmd/worker >> logs/worker.log &
nohup ./cmd/getData >> logs/getData.log &
