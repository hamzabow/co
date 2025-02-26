package prompts

var SimplePrompt = "Generate a concise and clear commit message describing " +
	"the following changes (output of `git diff`):\n```\n%s\n```\n\nEnsure the message is concise and meaningful. Return only the commit message, no extra text, and don't wrap the commit message with code blocks."

var ShortConventionalCommitsPrompt = "Generate a concise and clear commit message that follows" +
	" the **Conventional Commits** format. The commit message should" +
	" describe the following changes (output of command `git diff --staged`):" +
	"\n```\n%s\n```\n\n" +
	"Ensure the message is concise and meaningful. Return only the commit message," +
	" no extra text, and don't wrap the commit message with code blocks."

var LongConventionalCommitsPrompt = `Please generate a commit message following the **Conventional Commits** format.

---

### **Conventional Commits Specification**
A commit message consists of:
1. **A type**, which describes the category of change.
2. **An optional scope**, specifying the module or file affected.
3. **A concise description** of the change.
4. **An optional detailed body**, explaining what and why the change was made.
5. **Optional footers**, such as "BREAKING CHANGE" or issue references.

---

### **Format:**
` +
	"```" +
	`
<type>(<scope>): <short description>

[Optional body]

[Optional footers]
` +
	"```" +
	`

---

### **Types**
- **feat**: A new feature.
- **fix**: A bug fix.
- **docs**: Documentation changes only.
- **style**: Code style changes (formatting, missing semicolons, etc.) without affecting logic.
- **refactor**: A code change that neither fixes a bug nor adds a feature.
- **perf**: A performance improvement.
- **test**: Adding or updating tests.
- **build**: Changes to the build system or dependencies.
- **ci**: Changes to CI/CD pipelines.
- **chore**: Routine changes that do not affect production code.
- **revert**: A reversal of a previous commit.

---

### **Breaking Changes**
If a change introduces backward-incompatible behavior, it must include the ` + "`BREAKING CHANGE:`" + ` footer.

Example:
` + "```" + `
feat(api): update authentication method

The authentication API now requires a JWT token instead of a session ID.
Users must update their clients to use the new authentication flow.

BREAKING CHANGE: Authentication API now requires JWT tokens instead of session IDs.
` + "```" + `

---

### **Examples**
- "feat(auth): add OAuth2 login support"
- "fix(database): resolve race condition in transactions"
- "docs(readme): update installation instructions"
- "refactor(api): improve error handling"
- "perf(query): optimize search results loading time"
- "test(user): add unit tests for profile updates"
- "build(deps): upgrade TypeScript to v5.0.0"
- "ci(actions): update GitHub Actions workflow"
- "chore(lint): fix ESLint warnings"
- "revert(auth): rollback JWT token changes"

---

### **Given the following staged git diff, generate a commit message that strictly follows this specification:**

` + "```\n%s\n```\n\n" + `Ensure the message is concise and meaningful. Return only the commit message, no extra text, and don't wrap the commit message with code blocks.`

var GitmojiPrompt = "Generate a commit message that follows the **Gitmoji** " +
	"specification using the **Unicode format** for emojis.\n\n" +
	"### **Commit Message Format:**\n" +
	"- **Start with a Gitmoji in Unicode format** representing the intention of the change.\n" +
	"- **Optionally include a scope** in parentheses if relevant.\n" +
	"- **Use a colon (`:`) after the scope if provided**, or a space if not.\n" +
	"- **Write a brief and meaningful commit message**.\n\n" +
	"### **Format:**\n" +
	"<gitmoji> [optional(scope)][:?] <commit message>\n\n" +
	"### **Examples:**\n" +
	"⚡️ Lazyload home screen images\n" +
	"🐛 Fix `onClick` event handler\n" +
	"🔖 Bump version `1.2.0`\n" +
	"♻️ (components): Transform classes to hooks\n" +
	"📈 Add analytics to the dashboard\n" +
	"🌐 Support Japanese language\n" +
	"♿️ (account): Improve modals a11y\n\n" +
	"Ensure the commit message follows this format strictly. " +
	"Return only the commit message, no extra text, and don't wrap it with code blocks.\n\n" +
	"The commit message should describe the following changes (output of command `git diff --staged`):" +
	"\n```\n%s\n```"

var GitmojiShortcodePrompt = "Generate a commit message that follows the **Gitmoji** " +
	"specification using the **shortcode format** for emojis.\n\n" +
	"### **Commit Message Format:**\n" +
	"- **Start with a Gitmoji in Shortcode format** (e.g., `:zap:`, `:bug:`, `:bookmark:`).\n" +
	"- **Optionally include a scope** in parentheses if relevant.\n" +
	"- **Use a colon (`:`) after the scope if provided**, or a space if not.\n" +
	"- **Write a brief and meaningful commit message**.\n\n" +
	"### **Format:**\n" +
	":gitmoji: [optional(scope)][:?] <commit message>\n\n" +
	"### **Examples:**\n" +
	":zap: Lazyload home screen images\n" +
	":bug: Fix `onClick` event handler\n" +
	":bookmark: Bump version `1.2.0`\n" +
	":recycle: (components): Transform classes to hooks\n" +
	":chart_with_upwards_trend: Add analytics to the dashboard\n" +
	":globe_with_meridians: Support Japanese language\n" +
	":wheelchair: (account): Improve modals a11y\n\n" +
	"Ensure the commit message follows this format strictly. " +
	"Return only the commit message, no extra text, and don't wrap it with code blocks.\n\n" +
	"The commit message should describe the following changes (output of command `git diff --staged`):" +
	"\n```\n%s\n```"
