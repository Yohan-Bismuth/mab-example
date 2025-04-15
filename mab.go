package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Arm struct {
	Trials      int
	Successes   int
	TotalReward float64
	SuccessProb float64
	RewardValue float64
}

func (a *Arm) AvgReward() float64 {
	if a.Trials == 0 {
		return 0.0
	}
	return a.TotalReward / float64(a.Trials)
}

func (a *Arm) Play() {
	a.Trials++
	if rand.Float64() < a.SuccessProb {
		a.Successes++
		a.TotalReward += float64(rand.Intn(int(a.RewardValue)))
	}
}

func chooseArm(arms []Arm, epsilon float64) (int, bool) {
	avgRewards := []float64{arms[0].AvgReward(), arms[1].AvgReward()}

	exploring := rand.Float64() < epsilon
	if exploring {
		// Exploration : choisir celui qui semble moins bon
		if avgRewards[0] > avgRewards[1] {
			return 1, true
		} else if avgRewards[1] > avgRewards[0] {
			return 0, true
		}
		return rand.Intn(2), true
	}

	// Exploitation : choisir le meilleur bras (avec gestion égalité)
	if avgRewards[0] > avgRewards[1] {
		return 0, false
	} else if avgRewards[1] > avgRewards[0] {
		return 1, false
	}

	// En cas d'égalité, on choisit celui le moins joué
	if arms[0].Trials <= arms[1].Trials {
		return 0, false
	}
	return 1, false
}

func runEGSimulation() float64 {
	rand.Seed(time.Now().UnixNano())

	const (
		epsilon   = 0.15
		maxTrials = 30
	)

	arms := []Arm{
		{SuccessProb: 0.2, RewardValue: 20.0},
		{SuccessProb: 0.8, RewardValue: 2.0},
	}

	explorationCount := 0

	for i := 0; i < maxTrials; i++ {
		choice, exploring := chooseArm(arms, epsilon)
		if exploring {
			explorationCount++
		}
		arms[choice].Play()
	}

	/*
		for i, arm := range arms {
			fmt.Printf("Arm %d: %d successes / %d trials, total reward = %.2f, avg reward = %.2f\n",
				i, arm.Successes, arm.Trials, arm.TotalReward, arm.AvgReward())
		}
		fmt.Printf("Exploration count: %d\n", explorationCount)
	*/

	totalReward := 0.0
	for _, arm := range arms {
		totalReward += arm.AvgReward() * float64(arm.Trials)
	}
	return totalReward
}

func runUCBSimulation() float64 {
	rand.Seed(time.Now().UnixNano())
	const maxTrials = 30

	arms := [2]Arm{
		{SuccessProb: 0.2, RewardValue: 20.0},
		{SuccessProb: 0.8, RewardValue: 2.0},
	}

	// S'assurer que chaque bras est joué au moins une fois
	for i := range arms {
		arms[i].Play()
	}

	for t := len(arms); t < maxTrials; t++ {
		ucbValues := make([]float64, len(arms))
		totalTrials := t

		for i, arm := range arms {
			avg := arm.AvgReward()
			bonus := math.Sqrt((2 * math.Log(float64(totalTrials))) / float64(arm.Trials))
			ucbValues[i] = avg + bonus
		}

		// Choisir le bras avec la plus grande valeur UCB
		var choice int
		if ucbValues[0] > ucbValues[1] {
			choice = 0
		} else {
			choice = 1
		}
		// fmt.Printf("UCB values : [%f,%f]\n", ucbValues[0], ucbValues[1])

		arms[choice].Play()
	}

	/*for i, arm := range arms {
		fmt.Printf("Arm %d: %d successes / %d trials, total reward = %.2f, avg reward = %.2f\n",
			i, arm.Successes, arm.Trials, arm.TotalReward, arm.AvgReward())
	}*/

	return arms[0].TotalReward + arms[1].TotalReward
}

func main() {
	// Run EG simulation
	const nbSimulation = 1000
	total := 0.0
	for i := 0; i < nbSimulation; i++ {
		total += runEGSimulation()
	}
	fmt.Printf("Reward moyen epsilon greedy : %f\n", total/nbSimulation)

	// Run UCB simulation
	total = 0.0
	for i := 0; i < nbSimulation; i++ {
		total += runUCBSimulation()
	}
	fmt.Printf("Reward moyen UCB : %f\n", total/nbSimulation)
}
