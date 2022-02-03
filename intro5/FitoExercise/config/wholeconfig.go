package config

import(
  "encoding/json"
  "os"
  "fmt"
)

type WholeConfig struct {
  Dbusername string
  Dbpassword string
  Dbhostname string
  Dbschema   string
  Dedishost  string
  Serverbase string 
  Dbname     string
}

func Parse(path string) (ret *WholeConfig ,err error) {
  raw,err := os.ReadFile(path)
  if err != nil{
    return
  }
  err = json.Unmarshal(raw,&ret)
  return
}

func (config *WholeConfig) ConnectionString() string{
  return fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
                                   config.Dbhostname,
                                   config.Dbusername,
                                   config.Dbpassword,
                                   config.Dbname)
}
