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
    phi_eating []bool
}

func mayIEat(host Host, phi Philosopher) bool {
    /* Retorna True si el filosofo phi está habilitado para comer*/
    var people_eating = 0
    for i := 0; i < 5; i++ { // Revisa cuantas personas hay comiendo
        if host.phi_eating[i] {
            people_eating += 1
        }
    }
    if (people_eating < 2) { // Si hay menos de dos persoas intenta habilitar al filosofo
        AllowEating(host, phi)
        return true
    }
    return false
}

func startDinner(phi Philosopher, host Host, timesToEat int, wg *sync.WaitGroup) {
    /* 
    Se ejecuta cuando los filósofos comienzan la cena y corresponde a la función que permite que 
    estos coman 2 veces.
    */
    for i := 0; i < timesToEat; i++ { // timesToEat = 2
        ate := false
        for !ate { // Este for se ejecuta todo el tiempo hasta que ate=true
            if mayIEat(host, phi) { // Solicita al Host poder comer
                Eat(phi, host) 
                ate = true
            }
        }
    }
    defer wg.Done()
}

func Eat(phi Philosopher, host Host) {
    /*
    Corresponde a la función donde el filósofo come (ya tiene desbloqueados los palillos).
    Una vez que termina de comer desbloquea los palillos y se actualiza la información del host.
    */
	fmt.Printf("Filósofo %d comiendo\n", phi.num)
    time.Sleep(time.Second)
	fmt.Printf("Filósofo %d terminó de comer\n", phi.num)
    phi.leftHand.Unlock()
    phi.rightHand.Unlock()
    host.phi_eating[phi.num] = false
    time.Sleep(time.Second)
}

func AllowEating(host Host, phi Philosopher){
    /*
    Bloquea los palillos correspondientes al filosofo que quiere comer y actualiza la
    información del host.
    */
    phi.rightHand.Lock()
    phi.leftHand.Lock()
    host.phi_eating[phi.num] = true
}

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