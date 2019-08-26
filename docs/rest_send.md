

## rest  send


```
curl --location --request POST "http://127.0.0.1:1317/hs/send" \
  --header "Content-Type: application/x-www-form-urlencoded" \
  --data "    { \"base_req\": 
      { \"from\": \"htdf1l7spsv0lgx8npg3xvqfvsn8pdrpntv5djmmhuk\", 
        \"memo\": \"\",
        \"password\": \"12345678\", 
        \"chain_id\": \"testchain\", 
        \"account_number\": \"0\", 
        \"sequence\": \"0\", 
        \"gas\": \"200000\", 
        \"fees\": [ 
              { \"denom\": \"htdf\",
                 \"amount\": \"0.0000002\" } 
         ], 
         \"simulate\": false
      },          
      \"amount\": [ 
              { \"denom\": \"htdf\", 
                \"amount\": \"0.00001\" } ],
      \"to\": \"htdf1ucwpvw99u9tj3sxuwtnesge5tl90c6y0zcnc73\",
      \"data\": \"\"
    }"


```