rm -rf vendor
rm -rf Godeps
go get -u github.com/dafiti/echo-middleware
go get -u github.com/go-sql-driver/mysql
go get -u github.com/jinzhu/gorm
go get -u github.com/shopspring/decimal

go get -u gopkg.in/yaml.v2
# godep save
