export FILES_TO_EXCLUDE="'|$(yq e '.files' coverignore.yaml | tr '\n' '|' | tr -d '-' | tr -d [:blank:])'"


export DIRS_TO_EXCLUDE="'|$(yq e '.dirs' coverignore.yaml | tr '\n' '|' | tr -d '-' | tr -d [:blank:])'"

export PKGLIST=$(go list ./... | grep -v -E $DIRS_TO_EXCLUDE | tr '\n' ',')

ENV=test go test -coverprofile=cover.text.tmp ./internal/app/test/... -v -covermode=count -coverpkg=$PKGLIST -p 1 -ginkgo.noColor | go-junit-report > report.xml


cat cover.text.tmp | grep -v -E $FILES_TO_EXCLUDE > cover.out

go tool cover -html=cover.out -o cover.html
go tool cover -func cover.out

