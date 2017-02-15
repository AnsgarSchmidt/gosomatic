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
	LastModified         int64 `json:"lastModified"`
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
                 SpeedRatingBonus            float64 `json:"speedRatingBonus"`
                 Crit                        float64 `json:"crit"`
                 CritRating                  float64 `json:"critRating"`
                 Haste                       float64 `json:"haste"`
                 HasteRating                 float64 `json:"hasteRating"`
                 HasteRatingPercent          float64 `json:"hasteRatingPercent"`
                 Mastery                     float64 `json:"mastery"`
                 MasteryRating               float64 `json:"masteryRating"`
                 Leech                       float64 `json:"leech"`
                 LeechRating                 float64 `json:"leechRating"`
                 LeechRatingBonus            float64 `json:"leechRatingBonus"`
                 Versatility                 float64 `json:"versatility"`
                 VersatilityDamageDoneBonus  float64 `json:"versatilityDamageDoneBonus"`
                 VersatilityHealingDoneBonus float64 `json:"versatilityHealingDoneBonus"`
                 VersatilityDamageTakenBonus float64 `json:"versatilityDamageTakenBonus"`
                 AvoidanceRating             float64 `json:"avoidanceRating"`
                 AvoidanceRatingBonus        float64 `json:"avoidanceRatingBonus"`
                 SpellPen                    float64 `json:"spellPen"`
                 SpellCrit                   float64 `json:"spellCrit"`
                 SpellCritRating             float64 `json:"spellCritRating"`
                 Mana5                       float64 `json:"mana5"`
                 Mana5Combat                 float64 `json:"mana5Combat"`
                 Armor                       float64 `json:"armor"`
                 Dodge                       float64 `json:"dodge"`
                 DodgeRating                 float64 `json:"dodgeRating"`
                 Parry                       float64 `json:"parry"`
                 ParryRating                 float64 `json:"parryRating"`
                 Block                       float64 `json:"block"`
                 BlockRating                 float64 `json:"blockRating"`
                 MainHandDmgMin              float64 `json:"mainHandDmgMin"`
                 MainHandDmgMax              float64 `json:"mainHandDmgMax"`
                 MainHandSpeed               float64 `json:"mainHandSpeed"`
                 MainHandDps                 float64 `json:"mainHandDps"`
                 OffHandDmgMin               float64 `json:"offHandDmgMin"`
                 OffHandDmgMax               float64 `json:"offHandDmgMax"`
                 OffHandSpeed                float64 `json:"offHandSpeed"`
                 OffHandDps                  float64 `json:"offHandDps"`
                 RangedDmgMin                float64 `json:"rangedDmgMin"`
                 RangedDmgMax                float64 `json:"rangedDmgMax"`
                 RangedSpeed                 float64 `json:"rangedSpeed"`
                 RangedDps                   float64 `json:"rangedDps"`
	       }`json:"stats"`
}

func getWoWStatistic (url, key, realm, character string) (WoWStat, error) {
	var stat *WoWStat;
	req_url := fmt.Sprintf("%s%s/%s?fields=stats&locale=en_US&apikey=%s", url, realm, character, key)
	log.Println(req_url)
	resp, err := http.Get(req_url)
	if err != nil {
		log.Println("Error fetching data: %s for url: %s", err.Error(), url)
		return *stat, err
	}

	decoder := json.NewDecoder(resp.Body);
	if err = decoder.Decode(&stat); err != nil {
		log.Println("Error decoding feed: %s", err.Error())
		return *stat, err
	}
	return *stat, nil
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

	stat, err := getWoWStatistic(url, wowKey, realm, charName)

	if err == nil && stat.Level > 0 {
		text := fmt.Sprintf("%d", stat.Level)
		token := c.Publish("wow/Phawx/level", 0, false, text)
		token.Wait()

		text = fmt.Sprintf("%d", stat.Stats.Health)
		token = c.Publish("wow/Phawx/health", 0, false, text)
		token.Wait()

		text = fmt.Sprintf("%d", stat.Stats.Power)
		token = c.Publish("wow/Phawx/power", 0, false, text)
		token.Wait()

		text = fmt.Sprintf("%d", stat.TotalHonorableKills)
		token = c.Publish("wow/Phawx/totalHonorableKills", 0, false, text)
		token.Wait()
	}

	c.Disconnect(0)

}

