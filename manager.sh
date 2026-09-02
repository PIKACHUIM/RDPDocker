#!/bin/bash
function LIST_OCIS {
    source scripts/manager/titles.sh
    echo "   =============================List Container============================="
    sudo docker ps $1 --format '>  {{.ID}}   {{.Names}}   {{.Image}}   {{upper .State}}\n\t\t  Networks: {{.Networks | printf "%6s"}}  Status: {{.Status}}' 
    echo "   ==========================Enter Key to Continue========================="
}

while [ true ];
do
  source scripts/manager/titles.sh
  echo "   ==================================Menu=================================="
  echo "   --------------------------- Container Manage ---------------------------" 
  echo "          [n] [Created] Create a new docker container from an image"
  echo "          [l] [Display] Display all docker containers in the system"
  echo "          [s] [Started] Start your docker container already created"
  echo "          [t] [Pausing] Pause your docker container already running"
  echo "   -NOTE- [r] [Restart] Restart your docker container which ctarted"
  echo "   *WARN* [K] [Killall] Stop your docker container which is running"
  echo "   *WARN* [D] [Deleted] Select and delete your docker container now"
  echo "   ----------------------------- Image Manage -----------------------------" 
  echo "   -NOTE- [b] [ Build ] Build your docker image from the DockerFile"
  echo "   *WARN* [C] [ Clean ] Clean all unuse images and untag image file"
  echo "   *WARN* [P] [ Prune ] Prune system(all unuse images & containers)"
  echo "   ------------------------------------------------------------------------" 
  echo "          [q] [ Exit~ ] Do nothing, just exit this script! Byebye~~"
  echo "   ========================================================================"
  echo  
  OP_TYPE="."
  echo -n "   Choose Operation Type Number(*): "
  read OP_TYPE
  # Quit -----------------------------------------------------------------------------
  if [ "$OP_TYPE" = 'q' ]; then
    clear
    exit 0
  # creator ---------------------------------------------------------------------------
  elif [ $OP_TYPE == 'n' ]; then
    source scripts/manager/create.sh
  # List -----------------------------------------------------------------------------
  elif [ $OP_TYPE == 'l' ]; then
    LIST_OCIS -a
    read KEY
  # Start ----------------------------------------------------------------------------
  elif [ $OP_TYPE == 's' ]; then
    DEL_NAME=""
    LIST_OCIS -a
    echo -n "   Enter Container Name or ID to Launch(*): "
    read DEL_NAME
    docker start -- "$DEL_NAME"
  # Pause ----------------------------------------------------------------------------
  elif [ $OP_TYPE == 't' ]; then
    DEL_NAME=""
    LIST_OCIS " "
    echo -n "   Enter Container Name or ID to Stop  (*): "
    read DEL_NAME
    docker stop -- "$DEL_NAME"
  # Restart --------------------------------------------------------------------------
  elif [ $OP_TYPE == 'r' ]; then
    DEL_NAME=""
    LIST_OCIS " "
    echo -n "   Enter Container Name or ID to Stop  (*): "
    read DEL_NAME
    docker restart -- "$DEL_NAME"
  # Killall --------------------------------------------------------------------------
  elif [ $OP_TYPE == 'K' ]; then
    DEL_NAME=""
    LIST_OCIS " "
    echo -n "   Are You Sure Stop ALL Container Now(Y/N): "
    read DEL_NAME
    if [ $DEL_NAME == 'Y' ]; then
      docker stop $(docker ps -a -q)
    fi
  # Delete ---------------------------------------------------------------------------
  elif [ $OP_TYPE == 'D' ]; then
    DEL_NAME=""
    LIST_OCIS -a
    echo -n "   Enter Container Name or ID to Delete(*): "
    read DEL_NAME
    docker stop -- "$DEL_NAME"
    docker rm -- "$DEL_NAME"
  elif [ $OP_TYPE == 'b' ]; then
    source builder.sh
  elif [ $OP_TYPE == 'C' ]; then
    echo -n "   "
    docker rmi $(docker images -a | grep none | awk '{print $3}')
  elif [ $OP_TYPE == 'P' ]; then
    echo -n "   "
    docker system prune --volumes
  fi
done



