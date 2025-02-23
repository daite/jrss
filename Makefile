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

# 🏷️ Tag version, update main.go, and push to GitHub (tag should be in vX.X.X format with integers >= 0)
tag:
	@if [ -z "$(TAG)" ]; then \
		echo "❗ Usage: make tag TAG=vX.X.X"; \
		exit 1; \
	fi
	@if ! echo "$(TAG)" | grep -Eq '^v([0-9]+)\\.([0-9]+)\\.([0-9]+)$$'; then \
		echo "❌ Invalid tag format! Use vX.X.X format with integers greater than or equal to 0 (e.g., v1.0.0, v0.2.3, v2.1.5)"; \
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
	@echo "  make tag TAG=vX.X.X     - Update version in main.go, create an annotated Git tag (vX.X.X format, integers >= 0), and push to GitHub"