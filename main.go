package main

/*
* by akarizao - 2025
*/

import (
    "bufio"
    "bytes"
    "fmt"
    "os/exec"
    "strings"

    "github.com/manifoldco/promptui"
)

var authorMessage = "Autor: akarizao"

var translations = map[string]map[string]string{
    "pt": {
        "select_min":         "Selecione o valor mínimo (min_refresh_rate)",
        "select_peak":        "Selecione o valor máximo (peak_refresh_rate)",
        "configured":         "Configurado (min: %d, peak: %d)",
        "exit":               "Sair",
        "back":               "Voltar ao menu inicial",
        "reset":              "Redefinir taxas",
        "view_rates":         "Visualizar taxas atuais",
        "lang_select":        "Selecione o idioma",
        "lang_pt":            "Português",
        "lang_en":            "Inglês",
        "exit_message":       "Saindo...",
        "back_message":       "Voltando ao menu inicial...",
        "reset_message":      "Redefina as taxas de atualização.",
        "action_select":      "O que deseja fazer?",
        "current_rates":      "Taxas atuais: min_refresh_rate = %s, peak_refresh_rate = %s",
        "view_press_enter":   "Pressione Enter para voltar...",
        "error_reading_rates": "Erro ao ler taxas atuais",
        "set_rates_question": "Definir taxa de atualização?",
        "yes":                "Sim",
        "no":                 "Não",
    },
    "en": {
        "select_min":         "Select minimum value (min_refresh_rate)",
        "select_peak":        "Select maximum value (peak_refresh_rate)",
        "configured":         "Configured (min: %d, peak: %d)",
        "exit":               "Exit",
        "back":               "Back to main menu",
        "reset":              "Reset rates",
        "view_rates":         "View current rates",
        "lang_select":        "Select language",
        "lang_pt":            "Portuguese",
        "lang_en":            "English",
        "exit_message":       "Exiting...",
        "back_message":       "Returning to main menu...",
        "reset_message":      "Please reset the refresh rates.",
        "action_select":      "What would you like to do?",
        "current_rates":      "Current rates: min_refresh_rate = %s, peak_refresh_rate = %s",
        "view_press_enter":   "Press Enter to go back...",
        "error_reading_rates": "Error reading current rates",
        "set_rates_question": "Set refresh rate?",
        "yes":                "Yes",
        "no":                 "No",
    },
}

var refreshRates = []int{60, 90, 120}

func intToStringSlice(items []int) []string {
    s := make([]string, len(items))
    for i, v := range items {
        s[i] = fmt.Sprintf("%d Hz", v)
    }
    return s
}

func getSetting(setting string) (string, error) {
    cmd := exec.Command("settings", "get", "system", setting)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return "", err
    }
    result := strings.TrimSpace(out.String())
    return result, nil
}

func clearScreen() {
    fmt.Print("\033[H\033[2J")
}

func waitForEnter() {
    bufio.NewReaderSize(nil, 1).ReadBytes('\n')
}

func main() {
MainLoop:
    for {
        clearScreen()
        fmt.Println(authorMessage)
        languages := []string{"Português", "English"}
        codeMap := map[string]string{"Português": "pt", "English": "en"}

        langPrompt := promptui.Select{
            Label: translations["en"]["lang_select"],
            Items: languages,
        }
        _, langChosen, err := langPrompt.Run()
        if err != nil {
            fmt.Printf("Prompt failed: %v\n", err)
            return
        }

        lang := codeMap[langChosen]
        t := translations[lang]

        clearScreen()
        fmt.Println(authorMessage)
        setRatePrompt := promptui.Select{
            Label: t["set_rates_question"],
            Items: []string{t["yes"], t["no"]},
        }
        _, setRateChoice, err := setRatePrompt.Run()
        if err != nil {
            fmt.Printf("Prompt failed: %v\n", err)
            return
        }
        if setRateChoice == t["no"] {
            fmt.Println(t["exit_message"])
            return
        }

    RefreshLoop:
        for {
            clearScreen()
            fmt.Println(authorMessage)
            minOptions := intToStringSlice(refreshRates)
            minPrompt := promptui.Select{
                Label: t["select_min"],
                Items: minOptions,
            }
            _, minVal, err := minPrompt.Run()
            if err != nil {
                fmt.Printf("Prompt failed: %v\n", err)
                return
            }
            minRate := 0
            fmt.Sscanf(minVal, "%d Hz", &minRate)

            clearScreen()
            fmt.Println(authorMessage)
            peakPrompt := promptui.Select{
                Label: t["select_peak"],
                Items: minOptions,
            }
            _, peakVal, err := peakPrompt.Run()
            if err != nil {
                fmt.Printf("Prompt failed: %v\n", err)
                return
            }
            peakRate := 0
            fmt.Sscanf(peakVal, "%d Hz", &peakRate)

            exec.Command("settings", "put", "system", "min_refresh_rate", fmt.Sprintf("%d", minRate)).Run()
            exec.Command("settings", "put", "system", "peak_refresh_rate", fmt.Sprintf("%d", peakRate)).Run()

            fmt.Printf(t["configured"]+"\n", minRate, peakRate)

            for {
                fmt.Println(authorMessage)
                actionOptions := []string{t["back"], t["reset"], t["view_rates"], t["exit"]}
                actionPrompt := promptui.Select{
                    Label: t["action_select"],
                    Items: actionOptions,
                }
                _, action, err := actionPrompt.Run()
                if err != nil {
                    fmt.Printf("Prompt failed: %v\n", err)
                    return
                }

                switch action {
                case t["back"]:
                    fmt.Println(t["back_message"])
                    continue MainLoop
                case t["reset"]:
                    fmt.Println(t["reset_message"])
                    continue RefreshLoop
                case t["view_rates"]:
                    minCurrent, errMin := getSetting("min_refresh_rate")
                    peakCurrent, errPeak := getSetting("peak_refresh_rate")
                    if errMin != nil || errPeak != nil {
                        fmt.Println(t["error_reading_rates"])
                    } else {
                        fmt.Printf(t["current_rates"]+"\n", minCurrent, peakCurrent)
                    }
                    fmt.Println(t["view_press_enter"])
                    fmt.Scanln()
                case t["exit"]:
                    fmt.Println(t["exit_message"])
                    return
                }
            }
        }
    }
}
