tok=$(curl --header "Content-Type: application/json" --request POST\
    --data '{"username": "admin", "password": "password"}'\
    http://localhost:8080/login -q)
echo $tok
declare -a tokens
tokens[0]=$tok
function nc(){
    r=$1;
    i=$2;
    shift; shift;
    curl -q -H "Auth-Token: ${tokens[$i]}" \
        "$@"\
        http://localhost:8080/$r
}
function ncj(){
    a=$1
    b=$2
    c=$3
    shift;shift;shift
    nc $a $b --data "$c" "$@" -H "Content-Type: application/json" --request POST
}
nc token/test 0; echo

ncj post/create 0 '{"group_id": 1, "is_special_group": true, "body": "Hello World!"}'
nc post/getall 0

it=$(nc token/invite 0); echo $it
ncj token/reg 0 '{"username": "testuser", "password": "testuser"}' -H "Reg-Token: $it"
nt=$(ncj login 0 '{"username": "testuser", "password": "testuser"}');echo $nt
tokens[1]=$nt

ncj follow/follow 1 '{"username": "admin"}'

ncj post/getfollowing 1 '{"since": 0}'; echo

now=$(date +%s)
sleep 1
ncj post/create 0 '{"group_id": 2, "is_special_group": true, "body": "New Message"}'
sleep 5
ncj post/getfollowing 1 "{\"since\": $now}"; echo
sleep 2
ncj post/create 1 '{"group_id": 1, "is_special_group": true, "body": "Hello from User"}'
