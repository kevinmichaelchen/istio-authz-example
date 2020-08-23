clone_deps() {
  local provided=0
  if test -z "$1"; then
    echo "Arg 1 empty"
  else
    echo "Arg 1 not empty: $1"

    if test -z "$2"; then
      echo "Arg 2 empty"
    else
      echo "Arg 2 not empty: $2"
      provided=1
    fi
  fi

  if [ "$provided" -eq "0" ]; then
    # Using SSH keys
    echo "Using SSH"
    mkdir -p -m 0600 ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts && \
    git config --global url."git@github.com:".insteadOf "https://github.com/"
  else
    # Use access token
    echo "Using access token"

    printf "machine github.com\nlogin %s\npassword %s\n" "${GITHUB_USER}" "${GITHUB_ACCESS_TOKEN}" >> /root/.netrc
    chmod 600 /root/.netrc
  fi

  go mod download
}

main() {
  clone_deps "$@"
}

main "$@"
