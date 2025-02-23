# ⚙️ Variables
APP_NAME := jrss
MAIN_FILE := cmd/jrss/main.go
BUILD_DIR := bin

# Detect OS type for sed compatibility (Linux or macOS)
SED_INPLACE = $(shell if sed --version >/dev/null 2>&1; then echo "-i"; else echo "-i ''"; fi)

# 🏗️ Build the project
build:
	@echo "🔨 Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "✅ Build completed: $(BUILD_DIR)/$(APP_NAME)"

# 🚀 Run the project
run: build
	@echo "🏃 Running..."
	./$(BUILD_DIR)/$(APP_NAME)

# 🧹 Clean build artifacts
clean:
	@echo "🧹 Cleaning build files..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Clean completed!"

# 🏷️ Tag version and update main.go
tag:
	@if [ -z "$(TAG)" ]; then \
		echo "❗ Usage: make tag TAG=v1.0.0"; \
		exit 1; \
	fi
	@echo "🔄 Updating version in $(MAIN_FILE) to: $(TAG)"
	@sed $(SED_INPLACE) "s/^var Version = \".*\"/var Version = \"$(TAG)\"/" $(MAIN_FILE)
	git add $(MAIN_FILE)
	git commit -m "🔖 Version update: $(TAG)"
	git tag -a $(TAG) -m "🔖 Release $(TAG)"
	git push origin $(TAG)
	@echo "✅ Tag $(TAG) created and pushed to GitHub!"

# 🆘 Display help (default target)
help:
	@echo "📘 Available commands:"
	@echo "  make build              - Build the project"
	@echo "  make run                - Build and run the project"
	@echo "  make clean              - Remove build artifacts"
	@echo "  make tag TAG=v1.2.0     - Update version in main.go and create a Git tag"
