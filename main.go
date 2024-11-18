package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	tb "github.com/tucnak/telebot"
)

var resolvers = map[string]string{
	"Default":      "9.9.9.10",       // Default resolver
	"AdGuard":      "94.140.14.14",   // AdGuard
	"AT&T":         "165.87.13.129",  // AT&T
	"Cloudflare":   "1.1.1.1",        // Cloudflare
	"Comodo":       "8.26.56.26",     // Comodo
	"Google":       "8.8.8.8",        // Google
	"HiNet":        "168.95.1.1",     // HiNet
	"OpenDNS":      "208.67.222.222", // OpenDNS
	"Quad9":        "9.9.9.9",        // Quad9
	"Securolytics": "144.217.51.168", // Securolytics
	"UUNET-CH":     "195.129.12.122", // UUNET Switzerland
	"UUNET-DE":     "192.76.144.66",  // UUNET Germany
	"UUNET-UK":     "158.43.240.3",   // UUNET UK
	"UUNET-US":     "198.6.100.25",   // UUNET USA
	"Verisign":     "64.6.64.6",      // Verisign
	"Yandex":       "77.88.8.8",      // Yandex
}

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Warning: No .env file found.")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal().Msg("Error: BOT_TOKEN not set in environment.")
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create bot")
	}

	domainRegex := regexp.MustCompile(`^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	// Add /start command
	bot.Handle("/start", func(c *tb.Message) {
		startMessage := "Welcome to the DNS Resolver Bot! 🌐\n\n" +
			"I can help you resolve domain names using various DNS resolvers. Use the commands below to get started:\n" +
			"\nCommands:\n" +
			"/resolver - List available resolvers\n" +
			"`/lookup [resolver] domain` - Lookup a domain using a specific resolver\n" +
			"\nExample:\n" +
			"`/lookup Google` example.com\n" +
			"\nNeed help? Use /help."
		bot.Send(c.Sender, startMessage, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})

	// Add /help command
	bot.Handle("/help", func(c *tb.Message) {
		helpMessage := "Here's how to use this bot:\n\n" +
			"1️⃣ Use /resolver to see the available DNS resolvers.\n" +
			"2️⃣ Use /lookup [resolver] domain to resolve a domain using the specified resolver. If no resolver is specified, the default resolver is used.\n" +
			"\nExamples:\n" +
			"`/lookup Google` example.com\n" +
			"`/lookup example.com` (uses the default resolver)\n" +
			"\n🔧 Available Commands:\n" +
			"/start - Show welcome message\n" +
			"/help - Display this help message\n" +
			"/resolver - List available resolvers\n" +
			"/lookup - Resolve a domain\n\n" +
			"Happy resolving! 🚀"
		bot.Send(c.Sender, helpMessage, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})

	bot.Handle("/resolver", func(c *tb.Message) {
		resolverList := "Available resolvers:\n"
		for name := range resolvers {
			resolverList += fmt.Sprintf("- %s\n", name)
		}
		bot.Send(c.Sender, resolverList)
	})

	bot.Handle("/lookup", func(c *tb.Message) {
		args := strings.Fields(c.Payload)
		if len(args) == 0 {
			bot.Send(c.Sender, "Usage: /lookup [resolver] domain\nExample: /lookup Google example.com")
			return
		}

		var resolverName, domain string
		if len(args) == 1 {
			resolverName = "Default"
			domain = args[0]
		} else {
			resolverName = args[0]
			domain = args[1]
		}

		// Validate domain
		if !domainRegex.MatchString(domain) {
			bot.Send(c.Sender, "Error: Invalid domain format. Please use a valid domain like example.com.")
			return
		}

		// Get resolver IP
		resolverIP, ok := resolvers[resolverName]
		if !ok {
			bot.Send(c.Sender, "Error: Unknown resolver. Use /resolver to see the available resolvers.")
			return
		}

		// Perform DNS lookup
		ips, err := customDNSLookup(domain, resolverIP)
		if err != nil {
			bot.Send(c.Sender, fmt.Sprintf("Failed to resolve domain %s: %v", domain, err))
			return
		}

		// Send result
		result := fmt.Sprintf("Domain: %s\nResolver: %s (%s)\nIP Addresses: %s", domain, resolverName, resolverIP, strings.Join(ips, ", "))
		bot.Send(c.Sender, result)
	})

	// Handle inline queries
	bot.Handle(tb.OnQuery, func(q *tb.Query) {
		query := strings.TrimSpace(q.Text)
		parts := strings.SplitN(query, " ", 2)

		var resolverName, domain string
		if len(parts) == 1 {
			resolverName = "Default"
			domain = parts[0]
		} else {
			resolverName = parts[0]
			domain = parts[1]
		}

		// Validate domain
		if !domainRegex.MatchString(domain) {
			results := []tb.Result{
				&tb.ArticleResult{
					Title:       "Invalid Domain",
					Description: "Provide a valid domain like example.com",
					Text:        "Error: Invalid domain format.",
				},
			}
			bot.Answer(q, &tb.QueryResponse{Results: results})
			return
		}

		// Check resolver
		resolverIP, ok := resolvers[resolverName]
		if !ok {
			results := []tb.Result{
				&tb.ArticleResult{
					Title:       "Unknown Resolver",
					Description: "Available resolvers: Default, Google, Cloudflare, Quad9, OpenDNS, AdGuard",
					Text:        "Error: Unknown resolver. Use a valid resolver name.",
				},
			}
			bot.Answer(q, &tb.QueryResponse{Results: results})
			return
		}

		// Perform DNS lookup
		ips, err := customDNSLookup(domain, resolverIP)
		if err != nil {
			results := []tb.Result{
				&tb.ArticleResult{
					Title:       "DNS Lookup Failed",
					Description: fmt.Sprintf("Failed to resolve %s", domain),
					Text:        fmt.Sprintf("Failed to resolve %s: %v", domain, err),
				},
			}
			bot.Answer(q, &tb.QueryResponse{Results: results})
			return
		}

		// Inline response
		results := []tb.Result{
			&tb.ArticleResult{
				Title:       fmt.Sprintf("DNS Lookup for %s", domain),
				Description: fmt.Sprintf("Resolver: %s (%s)\nIPs: %s", resolverName, resolverIP, strings.Join(ips, ", ")),
				Text:        fmt.Sprintf("Domain: %s\nResolver: %s (%s)\nIPs: %s", domain, resolverName, resolverIP, strings.Join(ips, ", ")),
			},
		}
		bot.Answer(q, &tb.QueryResponse{Results: results})
	})

	// Start the bot
	log.Println("Bot is running...")
	bot.Start()
}

// Perform DNS lookup with a custom resolver
func customDNSLookup(domain, resolverIP string) ([]string, error) {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// Custom DNS resolver
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.DialContext(ctx, "udp", net.JoinHostPort(resolverIP, "53"))
		},
	}

	// Lookup IP addresses for the domain
	ips, err := resolver.LookupIPAddr(context.Background(), domain)
	if err != nil {
		return nil, err
	}

	// Convert to string slice
	ipStrings := make([]string, len(ips))
	for i, ip := range ips {
		ipStrings[i] = ip.IP.String()
	}
	return ipStrings, nil
}
