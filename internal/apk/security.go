package apk

import (
	"fmt"
)

// SecurityIssue represents a potential security issue found in the APK.
type SecurityIssue struct {
	Severity    string // HIGH, MEDIUM, LOW, INFO
	Category    string
	Description string
	Details     string
}

// DangerousPermission represents a dangerous Android permission.
type DangerousPermission struct {
	Name        string
	Description string
	Risk        string
}

// AbusivePermission represents a permission commonly abused by malware.
type AbusivePermission struct {
	Permission  string
	Status      string // normal, dangerous, unknown
	Info        string // short description
	Description string // detailed description
}

// AbusivePermissionResult represents the analysis result for abusive permissions.
type AbusivePermissionResult struct {
	TotalPermissions int
	MalwareMatches   []AbusivePermission
	OtherMatches     []AbusivePermission
	AllPermissions   []string
}

// AnalyzeAbusivePermissions checks APK permissions against known abusive permission lists.
func AnalyzeAbusivePermissions(apkPath string) (*AbusivePermissionResult, error) {
	info, err := ParseManifestBasic(apkPath)
	if err != nil {
		return nil, err
	}

	result := &AbusivePermissionResult{
		TotalPermissions: len(info.Permissions),
		AllPermissions:   info.Permissions,
	}

	for _, perm := range info.Permissions {
		if abusive, found := MalwarePermissions[perm]; found {
			result.MalwareMatches = append(result.MalwareMatches, abusive)
		}
		if abusive, found := OtherCommonAbusedPermissions[perm]; found {
			result.OtherMatches = append(result.OtherMatches, abusive)
		}
	}

	return result, nil
}

// AnalyzeSecurity performs security analysis on the APK.
// Based on OWASP MASTG security testing guidelines.
// Reference: https://mas.owasp.org/MASTG/
func AnalyzeSecurity(apkPath string) ([]SecurityIssue, error) {
	var issues []SecurityIssue

	info, err := GetAPKInfo(apkPath)
	if err != nil {
		return nil, err
	}

	// Check for debuggable flag (would need manifest parsing)
	if info.IsDebuggable {
		issues = append(issues, SecurityIssue{
			Severity:    "HIGH",
			Category:    "Configuration",
			Description: "Application is debuggable",
			Details:     "The android:debuggable flag is set to true. This allows attackers to attach debuggers and inspect app internals.",
		})
	}

	// Check for backup flag
	if info.AllowBackup {
		issues = append(issues, SecurityIssue{
			Severity:    "MEDIUM",
			Category:    "Configuration",
			Description: "Application allows backup",
			Details:     "The android:allowBackup flag is true. App data can be extracted via ADB backup.",
		})
	}

	// Search for sensitive patterns - Based on OWASP MASTG
	sensitivePatterns := map[string]SecurityIssue{
		// Hardcoded Secrets (MASTG-TEST-0001)
		`(?i)api[_-]?key\s*[=:]`:       {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential API key assignment found"},
		`(?i)api[_-]?secret`:           {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential API secret found"},
		`(?i)secret[_-]?key`:           {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential secret key found"},
		`(?i)private[_-]?key`:          {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential private key found"},
		`(?i)client[_-]?secret`:        {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential client secret found"},
		`(?i)auth[_-]?token`:           {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential auth token found"},
		`(?i)access[_-]?token`:         {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential access token found"},
		`(?i)bearer\s+[a-zA-Z0-9\-_]+`: {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential Bearer token found"},
		`(?i)password\s*[=:]`:          {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential hardcoded password found"},
		`(?i)passwd\s*[=:]`:            {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential hardcoded password found"},
		`(?i)encryption[_-]?key`:       {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential encryption key found"},
		`(?i)signing[_-]?key`:          {Severity: "HIGH", Category: "Hardcoded Secrets", Description: "Potential signing key found"},

		// Cloud Provider Credentials
		`AKIA[0-9A-Z]{16}`:        {Severity: "CRITICAL", Category: "Cloud Credentials", Description: "AWS Access Key ID found"},
		`(?i)aws[_-]?secret`:      {Severity: "CRITICAL", Category: "Cloud Credentials", Description: "Potential AWS secret key found"},
		`(?i)aws[_-]?access`:      {Severity: "HIGH", Category: "Cloud Credentials", Description: "Potential AWS access key found"},
		`(?i)aws[_-]?session`:     {Severity: "HIGH", Category: "Cloud Credentials", Description: "Potential AWS session token found"},
		`AIza[0-9A-Za-z\-_]{35}`:  {Severity: "HIGH", Category: "Cloud Credentials", Description: "Google API key found"},
		`(?i)gcp[_-]?api[_-]?key`: {Severity: "HIGH", Category: "Cloud Credentials", Description: "Potential GCP API key found"},
		`(?i)azure[_-]?`:          {Severity: "MEDIUM", Category: "Cloud Credentials", Description: "Azure reference found"},
		`(?i)digitalocean`:        {Severity: "MEDIUM", Category: "Cloud Credentials", Description: "DigitalOcean reference found"},
		`(?i)heroku`:              {Severity: "MEDIUM", Category: "Cloud Credentials", Description: "Heroku reference found"},

		// Firebase Configuration
		`(?i)firebase[_-]?api`:      {Severity: "MEDIUM", Category: "Firebase", Description: "Firebase API reference found"},
		`(?i)firebase[_-]?url`:      {Severity: "MEDIUM", Category: "Firebase", Description: "Firebase URL found"},
		`\.firebaseio\.com`:         {Severity: "MEDIUM", Category: "Firebase", Description: "Firebase Realtime Database URL found"},
		`\.firebaseapp\.com`:        {Severity: "INFO", Category: "Firebase", Description: "Firebase hosting URL found"},
		`(?i)firebase[_-]?database`: {Severity: "MEDIUM", Category: "Firebase", Description: "Firebase database reference found"},

		// Third-Party Services
		`(?i)stripe[_-]?`:             {Severity: "MEDIUM", Category: "Payment", Description: "Stripe payment integration detected"},
		`sk_live_[0-9a-zA-Z]{24}`:     {Severity: "CRITICAL", Category: "Payment", Description: "Stripe live secret key found"},
		`pk_live_[0-9a-zA-Z]{24}`:     {Severity: "HIGH", Category: "Payment", Description: "Stripe live publishable key found"},
		`(?i)paypal`:                  {Severity: "INFO", Category: "Payment", Description: "PayPal integration detected"},
		`(?i)braintree`:               {Severity: "INFO", Category: "Payment", Description: "Braintree integration detected"},
		`(?i)twilio`:                  {Severity: "MEDIUM", Category: "Third Party", Description: "Twilio integration detected"},
		`(?i)sendgrid`:                {Severity: "MEDIUM", Category: "Third Party", Description: "SendGrid integration detected"},
		`(?i)mailchimp`:               {Severity: "INFO", Category: "Third Party", Description: "Mailchimp integration detected"},
		`(?i)slack[_-]?webhook`:       {Severity: "MEDIUM", Category: "Third Party", Description: "Slack webhook found"},
		`xox[baprs]-[0-9a-zA-Z]{10,}`: {Severity: "HIGH", Category: "Third Party", Description: "Slack token found"},

		// Database Credentials
		`(?i)mongodb(\+srv)?://`: {Severity: "HIGH", Category: "Database", Description: "MongoDB connection string found"},
		`(?i)mysql://`:           {Severity: "HIGH", Category: "Database", Description: "MySQL connection string found"},
		`(?i)postgres(ql)?://`:   {Severity: "HIGH", Category: "Database", Description: "PostgreSQL connection string found"},
		`(?i)redis://`:           {Severity: "HIGH", Category: "Database", Description: "Redis connection string found"},
		`(?i)jdbc:`:              {Severity: "MEDIUM", Category: "Database", Description: "JDBC connection string found"},

		// Insecure Communication (MASTG-TEST-0006)
		`http://[^localhost][^\s"'<>]+`:            {Severity: "MEDIUM", Category: "Insecure Communication", Description: "HTTP URL found (non-HTTPS)"},
		`(?i)ssl[_-]?verify\s*[=:]\s*(false|0|no)`: {Severity: "HIGH", Category: "Insecure Communication", Description: "SSL verification disabled"},
		`(?i)trust[_-]?all[_-]?cert`:               {Severity: "HIGH", Category: "Insecure Communication", Description: "Trust all certificates pattern found"},
		`(?i)allow[_-]?all[_-]?hostname`:           {Severity: "HIGH", Category: "Insecure Communication", Description: "Allow all hostnames pattern found"},
		`(?i)insecure[_-]?ssl`:                     {Severity: "HIGH", Category: "Insecure Communication", Description: "Insecure SSL configuration found"},
		`setHostnameVerifier`:                      {Severity: "MEDIUM", Category: "Insecure Communication", Description: "Custom hostname verifier detected"},
		`TrustManager`:                             {Severity: "MEDIUM", Category: "Insecure Communication", Description: "Custom TrustManager detected"},
		`X509TrustManager`:                         {Severity: "MEDIUM", Category: "Insecure Communication", Description: "Custom X509TrustManager detected"},

		// Cryptography Issues (MASTG-TEST-0013, MASTG-TEST-0014)
		`(?i)DES[^C]`:          {Severity: "HIGH", Category: "Weak Cryptography", Description: "DES encryption (weak) detected"},
		`(?i)3DES`:             {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "3DES encryption (deprecated) detected"},
		`(?i)RC4`:              {Severity: "HIGH", Category: "Weak Cryptography", Description: "RC4 encryption (weak) detected"},
		`(?i)MD5`:              {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "MD5 hash (weak) detected"},
		`(?i)SHA-?1[^0-9]`:     {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "SHA-1 hash (deprecated) detected"},
		`(?i)ECB`:              {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "ECB mode (insecure) detected"},
		`(?i)PKCS1Padding`:     {Severity: "MEDIUM", Category: "Weak Cryptography", Description: "PKCS1Padding (vulnerable) detected"},
		`(?i)hardcoded[_-]?iv`: {Severity: "HIGH", Category: "Weak Cryptography", Description: "Hardcoded IV reference found"},
		`(?i)static[_-]?iv`:    {Severity: "HIGH", Category: "Weak Cryptography", Description: "Static IV reference found"},
		`SecureRandom`:         {Severity: "INFO", Category: "Cryptography", Description: "SecureRandom usage detected"},

		// Root/Jailbreak Detection
		`(?i)/system/app/Superuser`:        {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Superuser) found"},
		`(?i)/system/xbin/su`:              {Severity: "INFO", Category: "Root Detection", Description: "Root detection (su binary) found"},
		`(?i)/sbin/su`:                     {Severity: "INFO", Category: "Root Detection", Description: "Root detection (sbin su) found"},
		`(?i)com\.noshufou\.android\.su`:   {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Superuser app) found"},
		`(?i)com\.thirdparty\.superuser`:   {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Superuser) found"},
		`(?i)eu\.chainfire\.supersu`:       {Severity: "INFO", Category: "Root Detection", Description: "Root detection (SuperSU) found"},
		`(?i)com\.koushikdutta\.superuser`: {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Superuser) found"},
		`(?i)com\.topjohnwu\.magisk`:       {Severity: "INFO", Category: "Root Detection", Description: "Root detection (Magisk) found"},
		`(?i)isDeviceRooted`:               {Severity: "INFO", Category: "Root Detection", Description: "Root detection method found"},
		`(?i)checkRoot`:                    {Severity: "INFO", Category: "Root Detection", Description: "Root detection method found"},
		`(?i)RootBeer`:                     {Severity: "INFO", Category: "Root Detection", Description: "RootBeer library detected"},

		// Anti-Tampering & Anti-Debugging
		`(?i)frida`:               {Severity: "INFO", Category: "Anti-Tampering", Description: "Frida detection code found"},
		`(?i)xposed`:              {Severity: "INFO", Category: "Anti-Tampering", Description: "Xposed detection code found"},
		`(?i)substrate`:           {Severity: "INFO", Category: "Anti-Tampering", Description: "Cydia Substrate detection found"},
		`(?i)magisk`:              {Severity: "INFO", Category: "Anti-Tampering", Description: "Magisk detection code found"},
		`(?i)isDebuggerConnected`: {Severity: "INFO", Category: "Anti-Debugging", Description: "Debugger detection found"},
		`(?i)android\.os\.Debug`:  {Severity: "INFO", Category: "Anti-Debugging", Description: "Debug class usage found"},
		`(?i)ptrace`:              {Severity: "INFO", Category: "Anti-Debugging", Description: "Ptrace anti-debugging found"},
		`(?i)TracerPid`:           {Severity: "INFO", Category: "Anti-Debugging", Description: "TracerPid check found"},

		// Emulator Detection
		`(?i)generic.*Build`: {Severity: "INFO", Category: "Emulator Detection", Description: "Emulator detection pattern found"},
		`(?i)goldfish`:       {Severity: "INFO", Category: "Emulator Detection", Description: "Emulator detection (goldfish) found"},
		`(?i)isEmulator`:     {Severity: "INFO", Category: "Emulator Detection", Description: "Emulator detection method found"},

		// Logging & Debug Output
		`(?i)Log\.(d|v|i|w|e)\s*\(`: {Severity: "LOW", Category: "Information Disclosure", Description: "Android logging detected"},
		`(?i)System\.out\.print`:    {Severity: "LOW", Category: "Information Disclosure", Description: "System.out logging detected"},
		`(?i)printStackTrace`:       {Severity: "LOW", Category: "Information Disclosure", Description: "Stack trace printing detected"},
		`(?i)BuildConfig\.DEBUG`:    {Severity: "INFO", Category: "Debug Code", Description: "BuildConfig.DEBUG check found"},

		// WebView Issues (MASTG-TEST-0031)
		`setJavaScriptEnabled`:                {Severity: "MEDIUM", Category: "WebView", Description: "JavaScript enabled in WebView"},
		`addJavascriptInterface`:              {Severity: "HIGH", Category: "WebView", Description: "JavaScript interface in WebView (XSS risk)"},
		`setAllowFileAccess`:                  {Severity: "MEDIUM", Category: "WebView", Description: "File access in WebView"},
		`setAllowUniversalAccessFromFileURLs`: {Severity: "HIGH", Category: "WebView", Description: "Universal file access in WebView"},
		`setAllowFileAccessFromFileURLs`:      {Severity: "MEDIUM", Category: "WebView", Description: "File URL access in WebView"},

		// Data Storage (MASTG-TEST-0001, MASTG-TEST-0002)
		`MODE_WORLD_READABLE`:         {Severity: "HIGH", Category: "Data Storage", Description: "World-readable file mode found"},
		`MODE_WORLD_WRITEABLE`:        {Severity: "HIGH", Category: "Data Storage", Description: "World-writable file mode found"},
		`getExternalStorageDirectory`: {Severity: "MEDIUM", Category: "Data Storage", Description: "External storage usage detected"},
		`getExternalFilesDir`:         {Severity: "INFO", Category: "Data Storage", Description: "External files directory usage"},
		`SharedPreferences`:           {Severity: "INFO", Category: "Data Storage", Description: "SharedPreferences usage detected"},

		// Intent/IPC Vulnerabilities
		`(?i)intent\.setData`:       {Severity: "INFO", Category: "IPC", Description: "Intent data setting detected"},
		`(?i)intent\.putExtra`:      {Severity: "INFO", Category: "IPC", Description: "Intent extras detected"},
		`(?i)PendingIntent`:         {Severity: "INFO", Category: "IPC", Description: "PendingIntent usage detected"},
		`(?i)exported\s*=\s*"?true`: {Severity: "MEDIUM", Category: "IPC", Description: "Exported component detected"},

		// SQL Injection (MASTG-TEST-0026)
		`(?i)rawQuery`:            {Severity: "MEDIUM", Category: "SQL Injection", Description: "Raw SQL query detected"},
		`(?i)execSQL`:             {Severity: "MEDIUM", Category: "SQL Injection", Description: "execSQL usage detected"},
		`(?i)SELECT.*FROM.*WHERE`: {Severity: "LOW", Category: "SQL", Description: "SQL query pattern found"},

		// Backup Indicators
		`(?i)backup[_-]?agent`: {Severity: "MEDIUM", Category: "Backup", Description: "Backup agent reference found"},
		`(?i)BackupManager`:    {Severity: "INFO", Category: "Backup", Description: "BackupManager usage detected"},

		// Sensitive URLs/Endpoints
		`(?i)/api/`:                      {Severity: "INFO", Category: "API Endpoint", Description: "API endpoint path found"},
		`(?i)/v[0-9]+/`:                  {Severity: "INFO", Category: "API Endpoint", Description: "Versioned API path found"},
		`(?i)\.amazonaws\.com`:           {Severity: "MEDIUM", Category: "Cloud", Description: "AWS endpoint found"},
		`(?i)\.blob\.core\.windows\.net`: {Severity: "MEDIUM", Category: "Cloud", Description: "Azure blob storage found"},
		`(?i)\.s3\.`:                     {Severity: "MEDIUM", Category: "Cloud", Description: "AWS S3 bucket found"},

		// Biometric/Authentication
		`(?i)BiometricPrompt`:    {Severity: "INFO", Category: "Authentication", Description: "Biometric authentication detected"},
		`(?i)FingerprintManager`: {Severity: "INFO", Category: "Authentication", Description: "Fingerprint authentication detected"},
		`(?i)KeyguardManager`:    {Severity: "INFO", Category: "Authentication", Description: "Keyguard manager usage detected"},
	}

	for pattern, issue := range sensitivePatterns {
		matches, err := SearchStrings(apkPath, pattern)
		if err != nil {
			continue
		}
		if len(matches) > 0 {
			issue.Details = fmt.Sprintf("Found in %d file(s)", len(matches))
			issues = append(issues, issue)
		}
	}

	return issues, nil
}
