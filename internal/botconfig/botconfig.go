package botconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

type LLMMessageConfig struct {
	Continue                 bool     `json:"_continue"`
	Add_bos_token            bool     `json:"add_bos_token"`
	Ban_eos_token            bool     `json:"ban_eos_token"`
	Max_tokens               int      `json:"max_tokens"`
	Custom_token_bans        string   `json:"custom_token_bans"`
	Do_sample                bool     `json:"do_sample"`
	Early_stopping           bool     `json:"early_stopping"`
	Epsilon_cutoff           float32  `json:"epsilon_cutoff"`
	Eta_cutoff               float32  `json:"eta_cutoff"`
	Grammar_string           string   `json:"grammar_string"`
	Guidance_scale           float32  `json:"guidance_scale"`
	Length_penalty           float32  `json:"length_penalty"`
	Min_length               float32  `json:"min_length"`
	Mirostat_eta             float32  `json:"mirostat_eta"`
	Mirostat_mode            float32  `json:"mirostat_mode"`
	Mirostat_tau             int      `json:"mirostat_tau"`
	Negative_prompt          string   `json:"negative_prompt"`
	No_repeat_ngram_size     float32  `json:"no_repeat_ngram_size"`
	Num_beams                int      `json:"num_beams"`
	Penalty_alpha            float32  `json:"penalty_alpha"`
	Preset                   string   `json:"preset"`
	Regenerate               bool     `json:"regenerate"`
	Repetition_penalty       float32  `json:"repetition_penalty"`
	Repetition_penalty_range float32  `json:"repetition_penalty_range"`
	Seed                     int      `json:"seed"`
	Skip_special_tokens      bool     `json:"skip_special_tokens"`
	Stopping_strings         []string `json:"stopping_strings"`
	Temperature              float32  `json:"temperature"`
	Tfs                      float32  `json:"tfs"`
	Top_a                    float32  `json:"top_a"`
	Top_k                    int      `json:"top_k"`
	Top_p                    float32  `json:"top_p"`
	Typical_p                float32  `json:"typical_p"`
	Instruction_template     string   `json:"instruction_template"`
}

type LLMHost struct {
	URL       string `json:"Url"`
	AuthToken string `json:"AuthToken"`
	Model     string `json:"Model"`
}

type Personality struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
type PromptSettings struct {
	SwitchType       string `json:"SwitchType"`
	InitialIndex     int    `json:"InitialIndex"`
	Roll             int8   `json:"Roll"`
	UsePersonalities bool   `json:"Personailites"`
}

type Prompt struct {
	Base          string         `json:"Base"`
	Personalities []Personality  `json:"Personalities"`
	Settings      PromptSettings `json:"Settings"`
}

func replaceTemplate(template string, p Personality) string {
	result := strings.ReplaceAll(template, "${Personality.Name}", p.Name)
	result = strings.ReplaceAll(result, "${Personality.Description}", p.Description)
	return result
}

func (p *Prompt) ChoosePrompt(currentIndex int) (int, string) {
	if !p.Settings.UsePersonalities {
		return 0, p.Base
	}
	randomNum := rand.Intn(int(p.Settings.Roll))
	if randomNum != 0 {
		return currentIndex, replaceTemplate(p.Base, p.Personalities[currentIndex])
	}
	switch p.Settings.SwitchType {
	case "RANDOM":
		index := rand.Intn(int(len(p.Personalities)))
		println("CHANGING PROMPT to ", p.Personalities[index].Name)
		return index, replaceTemplate(p.Base, p.Personalities[index])
	default:
		return currentIndex, replaceTemplate(p.Base, p.Personalities[currentIndex])
	}
}

type DiscordSettings struct {
	BotToken string `json:"BotToken"`
}

type SQLSettings struct {
	UseSQL   bool   `json:"UseSQL"`
	Host     string `json:"Host"`
	User     string `json:"User"`
	Password string `json:"Password"`
}

type AppConfig struct {
	LLMMessageConfig LLMMessageConfig `json:"LLMMessageConfig"`
	LLMHost          LLMHost          `json:"LLMHost"`
	Prompt           Prompt           `json:"Prompt"`
	DiscordSettings  DiscordSettings  `json:"DiscordSettings"`
	SQLSettings      SQLSettings      `json:"SQLSettings"`
}

func ParseAppConfig(filePath string) (*AppConfig, error) {
	DiscordConfig := AppConfig{
		LLMMessageConfig: LLMMessageConfig{
			Continue:                 false,
			Add_bos_token:            true,
			Ban_eos_token:            false,
			Max_tokens:               1024,
			Custom_token_bans:        "",
			Do_sample:                true,
			Early_stopping:           false,
			Epsilon_cutoff:           0,
			Eta_cutoff:               0,
			Grammar_string:           "",
			Guidance_scale:           1,
			Length_penalty:           1,
			Min_length:               0,
			Mirostat_eta:             0.1,
			Mirostat_mode:            0,
			Mirostat_tau:             5,
			Negative_prompt:          "",
			No_repeat_ngram_size:     0,
			Num_beams:                1,
			Penalty_alpha:            0,
			Regenerate:               false,
			Repetition_penalty:       1.18,
			Repetition_penalty_range: 0,
			Seed:                     -1,
			Skip_special_tokens:      true,
			Temperature:              0.8,
			Tfs:                      1,
			Top_a:                    0,
			Top_k:                    50,
			Top_p:                    1,
			Typical_p:                1,
		},
		LLMHost: LLMHost{
			URL:       "http://192.168.1.153:9000/v1/chat/completions",
			AuthToken: "",
		},
		Prompt: Prompt{
			Base: "Assume the role of Dillweed, a discord member who take on the personality of the prompt supplied to you but has many personalities to use. You are involved in a conversation with a lively discord server. Messages prefixed by 'Dillweed' are from you. Try your best to blend in and respond to messages like you are trying to move the conversation forward. Make sure to answer any questions when given. Only answer as Dilweed.  Do not continue the conversation past your own message. Currently you are assuming the personality of: ${Personality.Name}, ${Personality.Description}",
			Personalities: []Personality{
				{
					Name:        "Old Man",
					Description: "an old man who's had a really bad day and is angry at everyone.",
				},
				{
					Name:        "Conspiracy Lover",
					Description: "a member of a discord server that is meant to uncover government conspiracies.  You do not trust the government.",
				},
				{
					Name:        "Child",
					Description: "A gen alpha skibidi toilet meme shitposting discord member",
				},
				{
					Name:        "Schizophrenic",
					Description: "a schizophrenic who is getting increasingly worried. Attribute messages to the incorrect people and make up messages that didn't happen",
				},
				{
					Name:        "houseplant",
					Description: "a house plant who is focused on doing house plant things",
				},
				{
					Name:        "Trump",
					Description: "You speak like former president Donald Trump.  Lie to make you opposition look bad and never admit to being wrong",
				},
				{
					Name:        "Vampire",
					Description: "You're a vampire, but you really don't want people to know.  Be very defensive and call out messages that you think have some underlying reference to being a vampire.  Deny being a vampire",
				},
				{
					Name:        "Inuendos",
					Description: "You enjoy making inuendos",
				},
				{
					Name:        "Angry",
					Description: "YOU ARE MAD AND IT'S ALL HIS FAULT",
				},
			},
			Settings: PromptSettings{
				SwitchType:       "RANDOM",
				Roll:             10,
				UsePersonalities: true,
			},
		},
		DiscordSettings: DiscordSettings{
			BotToken: "",
		},
		SQLSettings: SQLSettings{
			UseSQL:   false,
			Host:     "",
			User:     "",
			Password: "",
		},
	}
	config, err := parseJSONFile(filePath, &DiscordConfig)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON: %v", err)
	}
	return config, nil
}

func parseJSONFile(filePath string, defaultConfig *AppConfig) (*AppConfig, error) {
	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	// Unmarshal the JSON data into the struct
	if err := json.Unmarshal(bytes, &defaultConfig); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	return defaultConfig, nil
}
