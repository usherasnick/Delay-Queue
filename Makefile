PROJECT     := github.com/usherasnick/Delay-Queue
SRC         := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TARGETS     := Delay-Queue
ALL_TARGETS := $(TARGETS)

ifeq ($(race), 1)
	BUILD_FLAGS := -race
endif

ifeq ($(debug), 1)
	BUILD_FLAGS += -gcflags=all="-N -l"
endif

all: build

build: $(ALL_TARGETS)

$(TARGETS): $(SRC)
ifeq ("$(GOMODULEPATH)", "")
	@echo "no GOMODULEPATH env provided!!!"
	@exit 1
endif
	go build $(BUILD_FLAGS) $(GOMODULEPATH)/$(PROJECT)/cmd/$@

clean:
	rm -f $(ALL_TARGETS)

.PHONY: all build clean
