package repository

import (
	"database/sql"
        _ "github.com/lib/pq"
	"errors"
	"regexp"
)

type CryptoFilter struct {
  db * sql.DB
}

func  NewFilterRepository() (ret *CryptoFilter , err error){
   ret = & CryptoFilter{} 
   db,err:= sql.Open("postgres", 
        "host=ctolearn.cluster-ro-cjincqaxcmb8.us-east-1.rds.amazonaws.com "+
        "port=5432 user=postgres password=hdAXa4yVe7HWRXb dbname=postgres sslmode=disable")
        if err != nil{
          panic(err)
        }
  ret.db=db
  return 
}

func (f *CryptoFilter) List (from_datetime string, to_datetime string , 
                             use_mx bool, use_usd bool, use_hk bool) ( ret []map[string] string, err error){

  re:= regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}")
  if !re.Match([]byte (from_datetime)) || !re.Match([]byte(to_datetime)){
    err = errors.New("Invalid datetime format use yyyy-mm-dd hh:MM:ss")
    return
  }
  rows,err:=f.db.Query("SELECT mexican_peso::text, hk_dollar::text, usd_dollar::text FROM btc_loco.btc_loco"+
             " WHERE created_at >=$1 AND created_at <=$2",&from_datetime,&to_datetime) 
  if err != nil{
    return
  }
  for rows.Next(){
    n:=make(map[string] string)
    var mx,usd,hk string 
    rows.Scan(&mx,&usd,&hk)
    if use_mx {
      n["mx"]=mx
    }
    if use_usd{
      n["usd"]=usd
    }
    if use_hk{
      n["hk"]=hk
    }
    ret=append(ret,n)
  }
  return
}
