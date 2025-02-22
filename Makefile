.SILENT: table
table:
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
	cat stats/$$NAME.md;