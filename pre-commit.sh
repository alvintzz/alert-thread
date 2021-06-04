RED='\033[0;31m'
GREEN='\033[1;32m'
CYAN='\033[1;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

printUnitTest () {
	TOTAL_COVERAGE=0
	TOTAL_PACKAGE=0

	TEST_FLAG=false
	CURRENT_TEST=""
	IFS=$'\n'; for ROW in $1
	do
		if [[ $ROW == "?"* ]]
		then
			TOTAL_PACKAGE="$(($TOTAL_PACKAGE+1))"
			echo "${RED}$ROW${NC}"
			continue
		fi

		if [[ $ROW == "FAIL" ]]
		then
			continue
		fi
		
		if [[ $ROW == "FAIL"* ]]
		then
			TOTAL_PACKAGE="$(($TOTAL_PACKAGE+1))"
			echo "${RED}$ROW${NC}"
			continue
		fi

		if [[ $ROW == "PASS"* ]]
		then
			TOTAL_PACKAGE="$(($TOTAL_PACKAGE+1))"
			continue
		fi

		if [[ $ROW == "coverage: "* ]]
		then
			ROW=$(echo $ROW | sed 's/coverage: //')
			ROW=$(echo $ROW | sed 's/% of statements//')
			ROW=$(echo $ROW | sed 's/\.//')
			TOTAL_COVERAGE="$(($TOTAL_COVERAGE+$ROW))"
			continue
		fi

		if [[ $ROW == "ok"* ]]
		then
			echo "${GREEN}$ROW${NC}"
			continue
		fi

		if [[ $ROW == "=== RUN"* ]]
		then
			TEST_FLAG=true
			CURRENT_TEST="$ROW"
			continue
		fi
		if [[ $ROW == "--- PASS"* ]]
		then
			TEST_FLAG=false
			CURRENT_TEST=""
			continue
		fi
		if [[ $ROW == "--- FAIL"* ]]
		then
			echo "$CURRENT_TEST"
			TEST_FLAG=false
			CURRENT_TEST=""
			continue
		fi

		if [ $TEST_FLAG == false ]
		then
			echo $ROW
		fi
		CURRENT_TEST="$CURRENT_TEST\n$ROW"
	done

	AVERAGE_COVERAGE="$(($TOTAL_COVERAGE/$TOTAL_PACKAGE))"
	PERCENTAGE=$((AVERAGE_COVERAGE/10))
	LAST_DIGIT=${AVERAGE_COVERAGE: -1}

	echo "${YELLOW}=============================="
	echo "AVERAGE COVERAGE $PERCENTAGE.$LAST_DIGIT%"
	echo "==============================${NC}"
}

# Check whether the commit is success or not
RESULT=$?
if [ $RESULT -ne 0 ]
then
	echo "${RED}commit cancelled${NC}\n"
	exit 1
fi

# Do Unit Test to check the code validity and print the coverage information 
# Stop operation if the unit test is failed
echo "${CYAN}Testing before Commit${NC}"
UNIT_TEST=$(go test -v -race -cover $(go list ./... | grep -v /vendor/))
RESULT=$?
printUnitTest "$UNIT_TEST"

if [ $RESULT -ne 0 ]
then
	echo "${RED}Please fix above unit test issue${NC}\n"
	exit 1
fi

# Do basic linter for standarized code
echo "${CYAN}Lint the code${NC}"
golint ./...
LINTER=$?
if [ "$LINTER" -ne 0 ]
then
	echo "${RED}Please fix above basic linter issue${NC}"
	exit 1
fi

# Get all golang code that diff-ed
echo "${CYAN}Formatting golang files${NC}"
GOFILES=$(git diff --cached --name-only --diff-filter=ACM | grep '.go$')
if [ -z "$GOFILES" ]
then 
	echo "${GREEN}committed${NC}\n"
	exit 0
fi

# Get only the code that need to be formatted
UNFORMATTED=$(gofmt -l $GOFILES)
if [ -z "$UNFORMATTED" ]
then 
	echo "${GREEN}committed${NC}\n"
	exit 0
fi

# Format the code for each of golang file
for FILE in $UNFORMATTED; do
    echo "  gofmt -w ${GREEN}$PWD/$FILE${NC}"
    gofmt -w "$PWD/$FILE"
    git add "$PWD/$FILE"
done

echo "${GREEN}committed${NC}\n"
exit 0
