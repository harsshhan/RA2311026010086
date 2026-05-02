package scheduler

import (
	"github.com/harsshhan/RA2311026010086/vehicle_scheduling/models"
)

func Optimize(maxHours int, vehicles []models.Vehicle) (int, []string) {
	count := len(vehicles)
	
	memo := make([][]int, count+1)
	for i := range memo {
		memo[i] = make([]int, maxHours+1)
	}

	for i := 1; i <= count; i++ {
		veh := vehicles[i-1]
		
		for h := 0; h <= maxHours; h++ {
			if veh.Duration <= h {
				skip := memo[i-1][h]
				take := memo[i-1][h-veh.Duration] + veh.Impact
				
				if skip > take {
					memo[i][h] = skip
				} else {
					memo[i][h] = take
				}
			} else {
				memo[i][h] = memo[i-1][h]
			}
		}
	}

	var selected []string
	impact := memo[count][maxHours]
	hours := maxHours

	for i := count; i > 0 && impact > 0; i-- {
		if impact != memo[i-1][hours] {
			veh := vehicles[i-1]
			selected = append(selected, veh.TaskID)
			
			impact -= veh.Impact
			hours -= veh.Duration
		}
	}

	return memo[count][maxHours], selected
}
