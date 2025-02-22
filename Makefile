all: portifolio sums

.SILENT: portifolio
portifolio:
	NAME=`date +"%Y-%m-%d"` && \
	echo "| Model | Ticket | Sum | Quantity |" > stats/$$NAME.md && \
	echo "|-------|-------|-------|--------|" >> stats/$$NAME.md && \
	for model in `ls orders`; \
		do go run ./cmd/list --model "$$model" --roundsums \
		| jq -r ".holdings[] | [\"$$model\", .ticket, .sum, .quantity] \
		| @csv" \
		| tr ',' '|' \
		| tr '"' '`' \
		| awk '{print "|"$$0"|"}' >> stats/$$NAME.md; \
	done; \
	cat stats/$$NAME.md && echo -e "\n"

.SILENT: sums
sums:
	NAME=`date +"%Y-%m-%d"` && \
	echo "| Model | Total Sum | Change |" > stats/$$NAME-sums.md && \
	echo "|-------|-----------|--------|" >> stats/$$NAME-sums.md && \
	for model in `ls orders`; \
		do go run ./cmd/list --model "$$model" --roundsums \
		| jq -r '.holdings[] | .sum' \
		| awk '{sum += $$1} END {print "'"$$model"'" "\t" sum}' \
		>> stats/$$NAME-sums.temp; \
	done; \
	sort -k2 -nr stats/$$NAME-sums.temp \
	| awk '{print "|`" $$1 "`|" $$2 "|" "—" "|"}' \
	>> stats/$$NAME-sums.md; \
	rm stats/$$NAME-sums.temp; \
	cat stats/$$NAME-sums.md; \
	rm stats/$$NAME-sums.md;
