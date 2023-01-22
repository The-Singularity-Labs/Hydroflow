
#
# Example City Pipeline
# Author: Hippodamus
#
TARGET_DIR?=/tmp

ifeq ($(OS),Windows_NT)
	uname_S := Windows
else
	uname_S := $(shell uname -s)
endif

all: $(TARGET_DIR)/show.md

/tmp:
	mkdir -p /tmp


$(TARGET_DIR)/cities.csv: $(TARGET_DIR)  
	wget -O /tmp/cities.csv https://raw.githubusercontent.com/icyrockcom/country-capitals/master/data/country-list.csv

$(TARGET_DIR)/cities.db: $(TARGET_DIR) $(TARGET_DIR)/cities.csv  
	rm -rf /tmp/cities.db && sqlite3 /tmp/cities.db -cmd '.mode csv' '.import /tmp/cities.csv cities'

$(TARGET_DIR)/calculate_count.md: $(TARGET_DIR) $(TARGET_DIR)/cities.db  
	out=$$(sqlite3 /tmp/cities.db -header -csv 'SELECT COUNT(*) FROM cities') && echo "$$out"  > /tmp/calculate_count.md

$(TARGET_DIR)/preview.md: $(TARGET_DIR) $(TARGET_DIR)/cities.db  
	out=$$(sqlite3 /tmp/cities.db -header -csv 'SELECT * FROM cities LIMIT 5') && echo "$$out" > /tmp/preview.md

$(TARGET_DIR)/show.md: $(TARGET_DIR) $(TARGET_DIR)/calculate_count.md $(TARGET_DIR)/preview.md  
	cat /tmp/calculate_count.md
	cat /tmp/preview.md
	touch /tmp/show.md
	

