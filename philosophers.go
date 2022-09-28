package main

import (
    "fmt"
    "sync"
    "time"
)

	
type Chopstick struct {
    sync.Mutex
    id_ int
}

type Philosopher struct {
    num int
    leftHand *Chopstick
    rightHand *Chopstick
}

type Host struct {
    phy_eating []bool
}

func mayIEat(host Host, phy Philosopher) bool {
    var people_eating = 0
    for i := 0; i < 5; i++ {
        if host.phy_eating[i] {
            people_eating += 1
        }
    }
    if (people_eating < 2) {
        AllowEating(host, phy)
        return true
    }
    return false
}

func startDinner(phy Philosopher, host Host, timesToEat int, wg *sync.WaitGroup) {
    for i := 0; i < timesToEat; i++ {
        ate := false
        for !ate {
            if mayIEat(host, phy) {
                Eat(phy, host)
                ate = true
            }
        }

    }
    defer wg.Done()
}

func Eat(phy Philosopher, host Host) {
	fmt.Printf("Filósofo %d comiendo\n", phy.num)
    time.Sleep(time.Second)
	fmt.Printf("Filósofo %d terminó de comer\n", phy.num)
    phy.leftHand.Unlock()
    phy.rightHand.Unlock()
    host.phy_eating[phy.num] = false
    time.Sleep(time.Second)
}

func AllowEating(host Host, phy Philosopher){
    /*
    DEBE SER ALEATORIO EL PALILLO QUE SE BLOQUEA PRIMERO !!!!
    ESPERAR RESPUESTA AYUDANTE SOBRE SI SE PUEDE USAR RANDOM...
    POR MIENTRAS SE SELECCIONARÁ PRIMERO EL DERECHO
    */
    phy.rightHand.Lock()
    phy.leftHand.Lock()
    host.phy_eating[phy.num] = true
}

// Channels: y WaitGrroups


func main() { 
    var wg sync.WaitGroup
    var chopsticks = []Chopstick{Chopstick{}, Chopstick{}, Chopstick{}, Chopstick{}, Chopstick{}}
    for i := 0; i < 5; i++ {
        chopsticks[i].id_ = i
    }
    var philosophers = []Philosopher{Philosopher{0, &chopsticks[4], &chopsticks[0]}, Philosopher{1, &chopsticks[0], &chopsticks[1]}, Philosopher{2, &chopsticks[1], &chopsticks[2]}, Philosopher{3, &chopsticks[2], &chopsticks[3]}, Philosopher{4, &chopsticks[3], &chopsticks[4]}}
    var host = Host{[]bool {false, false, false, false, false} }
    const timesToEat = 2 // Veces que come cada filósofo

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go startDinner(philosophers[i], host, timesToEat, &wg)
    }
    wg.Wait() // Esperamos que todos los filósofos terminen de comer
}