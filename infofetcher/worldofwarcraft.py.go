package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"os"
	"time"
	"log"
	"encoding/json"
	"net/http"
)

type WoWStat struct {
	LastModified           int `json:"lastModified"`
	Name                string `json:"name"`
	Realm               string `json:"realm"`
        BattleGroup         string `json:"battlegroup"`
        Class                  int `json:"class"`
        Race                   int `json:"race"`
        Gender                 int `json:"gender"`
        Level                  int `json:"level"`
        AchievementsPoints     int `json:"achievementPoints"`
        Thumbnail           string `json:"thumbnail"`
        CalcClass           string `json:"calcClass"`
        Faction                int `json:"faction"`
	TotalHonorableKills    int `json:"totalHonorableKills"`
        Stats  struct{
		 Health                          int `json:"health"`
		 PowerType                    string `json:"powerType"`
                 Power                           int `json:"power"`
                 Str                             int `json:"str"`
                 Agi                             int `json:"agi"`
                 Int                             int `json:"int"`
                 Sta                             int `json:"sta"`
                 SpeedRating                 float64 `json:"speedRating"`
                 SpeedRatingBonus                int `json:"speedRatingBonus"`
                 Crit                        float64 `json:"crit"`
                 CritRating                      int `json:"critRating"`
                 Haste                       float64 `json:"haste"`
                 HasteRating                     int `json:"hasteRating"`
                 HasteRatingPercent          float64 `json:"hasteRatingPercent"`
                 Mastery                     float64 `json:"mastery"`
                 MasteryRating                   int `json:"masteryRating"`
                 Leech                           int `json:"leech"`
                 LeechRating                     int `json:"leechRating"`
                 LeechRatingBonus                int `json:"leechRatingBonus"`
                 Versatility                     int `json:"versatility"`
                 VersatilityDamageDoneBonus  float64 `json:"versatilityDamageDoneBonus"`
                 VersatilityHealingDoneBonus float64 `json:"versatilityHealingDoneBonus"`
                 VersatilityDamageTakenBonus float64 `json:"versatilityDamageTakenBonus"`
                 AvoidanceRating                 int `json:"avoidanceRating"`
                 AvoidanceRatingBonus        float64 `json:"avoidanceRatingBonus"`
                 SpellPen                        int `json:"spellPen"`
                 SpellCrit                   float64 `json:"spellCrit"`
                 SpellCritRating                 int `json:"spellCritRating"`
                 Mana5                           int `json:"mana5"`
                 Mana5Combat                     int `json:"mana5Combat"`
                 Armor                           int `json:"armor"`
                 Dodge                           int `json:"dodge"`
                 DodgeRating                     int `json:"dodgeRating"`
                 Parry                           int `json:"parry"`
                 ParryRating                     int `json:"parryRating"`
                 Block                           int `json:"block"`
                 BlockRating                     int `json:"blockRating"`
                 MainHandDmgMin                  int `json:"mainHandDmgMin"`
                 MainHandDmgMax                  int `json:"mainHandDmgMax"`
                 MainHandSpeed               float64 `json:"mainHandSpeed"`
                 MainHandDps                 float64 `json:"mainHandDps"`
                 OffHandDmgMin                   int `json:"offHandDmgMin"`
                 OffHandDmgMax                   int `json:"offHandDmgMax"`
                 OffHandSpeed                float64 `json:"offHandSpeed"`
                 OffHandDps                  float64 `json:"offHandDps"`
                 RangedDmgMin                    int `json:"rangedDmgMin"`
                 RangedDmgMax                    int `json:"rangedDmgMax"`
                 RangedSpeed                     int `json:"rangedSpeed"`
                 RangedDps                       int `json:"rangedDps"`
	       }`json:"stats"`
}

func getWoWStatistic (url, key, realm, character string) (WoWStat) {
	var stat *WoWStat;
	req_url := fmt.Sprintf("%s%s/%s?fields=stats&locale=en_US&apikey=%s", url, realm, character, key)
	log.Println(req_url)
	resp, err := http.Get(req_url)
	if err != nil {
		log.Println("Error fetching data: %s for url: %s", err.Error(), url)
		return *stat
	}

	decoder := json.NewDecoder(resp.Body);
	if err = decoder.Decode(&stat); err != nil {
		log.Println("Error decoding feed: %s", err.Error())
		return *stat
	}
	return *stat
}

func main() {

	fmt.Println("Word of Warcraft updater")

	url        := "https://us.api.battle.net/wow/character/"
	realm      := os.Getenv("REALM")
	charName   := os.Getenv("CHARACTER")
	wowKey     := os.Getenv("WOW_API_KEY")
	mqttServer := os.Getenv("MQTT_SERVER")

	if realm == "" {
		fmt.Print("REALM not set")
		os.Exit(1)
	}

	if charName == "" {
		fmt.Print("CHARACTER not set")
		os.Exit(1)
	}

	if wowKey == "" {
		fmt.Print("WOW_API_KEY not set")
		os.Exit(1)
	}

	if mqttServer == "" {
		fmt.Print("MQTT_SERVER not set")
		os.Exit(1)
	}

	// MQTT
	opts := mqtt.NewClientOptions().AddBroker(mqttServer).SetClientID("gowordofwarcraft")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
		os.Exit(2)
	}

	stat := getWoWStatistic(url, wowKey, realm, charName)

	text  := fmt.Sprintf("%d", stat.Level)
	token := c.Publish("wow/Phawx/level", 0, false, text)
	token.Wait()

	text  = fmt.Sprintf("%d", stat.Stats.Health)
	token = c.Publish("wow/Phawx/health", 0, false, text)
	token.Wait()

	text  = fmt.Sprintf("%d", stat.Stats.Power)
	token = c.Publish("wow/Phawx/power", 0, false, text)
	token.Wait()

	text  = fmt.Sprintf("%d", stat.TotalHonorableKills)
	token = c.Publish("wow/Phawx/totalHonorableKills", 0, false, text)
	token.Wait()

	c.Disconnect(0)

}

