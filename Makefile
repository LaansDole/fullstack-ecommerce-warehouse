SHELL := /bin/bash
.PHONY: publish
publish:
	git add .
	@while true; do \
		echo "Enter commit type:"; \
		echo "1. feature"; \
		echo "2. fix"; \
		echo "3. update"; \
		echo "4. fix & update"; \
		read -p "Choose a number from 1-4: " num; \
		case $$num in \
			1) type="feature"; break;; \
			2) type="fix"; break;; \
			3) type="update"; break;; \
			4) type="fix & update"; break;; \
			*) echo "Please choose again";; \
		esac; \
	done; \
	read -p "Enter commit message: " message; \
	git commit -m "$$type: $$message"
	git fetch
	git pull
	git push -u