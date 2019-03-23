if [ -z $1 ]; then
	URL=http://localhost:8080
else
	URL=$1
fi
DOCHK=1

chk(){
    return $DOCHK
}
msg(){
    if  chk;  then return; fi
    echo; echo $@;
    sleep 1
    read -p "> " yn
}   
msg Logging in as admin
tok=$(curl --header "Content-Type: application/json" --request POST\
    --data '{"username": "admin", "password": "password"}'\
    $URL/login -q)
echo $tok
declare -a tokens
tokens[0]=$tok
function nc(){
    r=$1;
    i=$2;
    shift; shift;
    curl -q -H "Auth-Token: ${tokens[$i]}" \
        "$@"\
        $URL/$r
    exit $?
}
function ncj(){
    a=$1
    b=$2
    c=$3
    shift;shift;shift
    nc $a $b --data "$c" "$@" -H "Content-Type: application/json" --request POST
}

msg Checking token
nc token/test 0; echo

msg Creating a post
ncj post/create 0 '{"group_id": 1, "is_special_group": true, "body": "Hello World!"}'

msg Fetching all posts
nc post/getall 0

msg Creating an invite token
it=$(nc token/invite 0); echo $it
msg Creating a new user
ncj token/reg 0 '{"username": "testuser", "password": "testuser"}' -H "Reg-Token: $it"
msg Logging in as new user
nt=$(ncj login 0 '{"username": "testuser", "password": "testuser"}');echo $nt
tokens[1]=$nt

msg Following admin as new user
ncj follow/follow 1 '{"username": "admin"}'

msg Get all posts from followed users as new user
ncj post/getfollowing 1 '{"since": 0}'; echo

now=$(date +%s)
msg Creating a new post as admin
ncj post/create 0 '{"group_id": 2, "is_special_group": true, "body": "New Message"}'
sleep 1
msg Fetching the new post as the new user using a timestamp
ncj post/getfollowing 1 "{\"since\": $now}"; echo
msg Creating a post as the new user
ncj post/create 1 '{"group_id": 1, "is_special_group": true, "body": "Hello from User"}'
msg Showing users the new user follows
nc follow/getfollowed 1 
msg Showing users following admin
nc follow/getfollowers 0
msg Creating a new group
gid=$(ncj group/create 0 '{"name":"generic"}'); echo $gid
