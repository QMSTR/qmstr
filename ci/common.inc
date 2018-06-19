function create_dirs() {
    if [ -n ${USERID} ]; then
        for i in /var/cache/qmstr /var/qmstr /var/lib/qmstr; do
            echo "Creating dir $i to be owned by ${USERID}"
            mkdir ${i}
            chown ${USERID} ${i}
        done
    fi
}

function fix_permission() {
    if [ -n ${USERID} ]; then
        for i in /usr/local/lib/python2.7/dist-packages/extractcode; do
            echo "Fixing $i to be owned by ${USERID}"
            chown -R ${USERID} ${i}
        done
    fi
}

function create_qmstr_user() {
    if [ -n ${USERID} ]; then
        useradd -m -u ${USERID} qmstr
        create_dirs
        fix_permission
    fi
}

function get_qmstr_user() {
    if [ -n ${USERID} ]; then
        getent passwd ${USERID} | cut -d ":" -f 1
    fi
}

function start_dgraph() {
    dgraph version
    dgraph zero &
    dgraph server --lru_mb=2048 --zero=localhost:5080 &
}

function start_dgraph_web() {
    dgraph-ratel &
}

QMSTR_CMDLINE="/usr/local/bin/qmstr-master --config /buildroot/qmstr.yaml ${PATH_SUB:+--pathsub=\"$PATH_SUB\"}"

function start_qmstr() {
    echo "$(get_qmstr_user) starting qmstr-master"
    su -c "${QMSTR_CMDLINE}" $(get_qmstr_user)
}
