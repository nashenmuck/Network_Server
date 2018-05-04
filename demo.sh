tok=$(curl --header "Content-Type: application/json" --request POST\
    --data '{"username": "admin", "password": "password"}'\
    http://localhost:8080/login -q)
echo $tok
function nc(){
    r=$1;
    shift
    curl -q -H "Auth-Token: $tok" \
        $@\
        http://localhost:8080/$r
}
function ncj(){
    r=$1;
     curl -q -H "Auth-Token: $tok" \
        -H "Content-Type: application/json"\
        --request POST\
        --data "$2"\
        http://localhost:8080/$r
}
nc token/test; echo

ncj post/create '{"target": 1, "special": true, "post": "Hello World!"}'

nc post/getall
