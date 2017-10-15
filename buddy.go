package main

import (
    "encoding/csv"
    "flag"
    "fmt"
    "github.com/golang/glog"
    "math/rand"
    "os"
    "time"
)

type worker struct {
    email string
    team string
}

func indices(head []string) (email int, team int, err error) {
    email = -1
    team = -1
    for i, col := range head {
        switch col {
            case "email":
                email = i
            case "team":
                team = i
        }
    }
    if email == -1 || team == -1 {
        err = fmt.Errorf("could not get all columns")
    }
    return
}

func read(in *string) ([]worker, error) {
    glog.V(2).Infof("opening file: %v", *in)
    f, err := os.Open(*in)
    if err != nil {
        glog.V(2).Infof("unable to open file: %v", *in)
        return nil, fmt.Errorf("could not open file: %v", err)
    }
    r := csv.NewReader(f)
    records, err := r.ReadAll()
    if err != nil {
        return nil, err
    }
    if len(records) < 2 {
        return nil, fmt.Errorf("not enough rows: %v", records)
    }
    email, team, err := indices(records[0])
    if err != nil {
        return nil, fmt.Errorf("did not have appropriate columns %v, err: %v", records[0], err)
    }
    var workers []worker
    for _, rec := range records[1:] {
        w := worker{email: rec[email], team: rec[team]}
        glog.V(2).Infof("got worker %v", w)
        workers = append(workers, w)
    }
    return workers, nil
}

func groupByTeam(workers []worker) [][]worker {
    teamNums := make(map[string]int)
    nextInt := 0
    teams := [][]worker{}
    for _, w := range workers {
        tNum, prs := teamNums[w.team]
        if !prs {
            teamNums[w.team] = nextInt
            tNum = nextInt
            glog.V(2).Infof("new team %v, adding worker %v", tNum, w)
            teams = append(teams, []worker{w})
            nextInt++
        } else {
            glog.V(2).Infof("adding worker %v to team %v", w, tNum)
            teams[tNum] = append(teams[tNum], w)
        }
    }
    return teams
}

func teams2group(teamCount []int, numGroups int) [][]int {
    groups := make([][]int, numGroups)
    lastGroup := 0
    for i, c := range teamCount {
        for n := 0; n < c; n++ {
            groups[lastGroup] = append(groups[lastGroup], i)
            lastGroup++
            lastGroup = lastGroup % numGroups
        }
    }
    return groups
}

func makeBuddies(workers []worker, size *int) [][]string {
    numGroups := len(workers) / *size
    teams := groupByTeam(workers)
    teamCounts := []int{}
    for _, t := range teams {
        glog.V(2).Infof("team %v has %v workers", t, len(t))
        teamCounts = append(teamCounts, len(t))
    }
    grouped := teams2group(teamCounts, numGroups)
    glog.V(2).Infof("grouped: %v", grouped)
    buddies := [][]string{}
    for _, b := range grouped {
        bGrp := []string{}
        for _, t := range b {
            glog.V(2).Infof("getting worker for team %v (%v) and group %v", t, teams[t], b)
            ind := 0
            l := len(teams[t]) - 1
            if l > 0 {
                ind = rand.Intn(l)
            }
            w := teams[t][ind]
            glog.V(2).Infof("got worker %v for group %v", w, b)
            teams[t] = append(teams[t][:ind], teams[t][ind+1:]...)
            bGrp = append(bGrp, w.email)
        }
        glog.V(2).Infof("added group %v", bGrp)
        buddies = append(buddies, bGrp)
    }
    return buddies
}

func main() {
    inp := flag.String("input_file", "./tmp/workers.csv", "A csv of emails and teams to import")
    rndm := flag.Bool("randomize", true, "A csv of emails and teams to import")
    grpSz := flag.Int("group_size", 4, "The number of people in a group")
    rand.Seed(time.Now().UTC().UnixNano())
    flag.Parse()
    workers, err := read(inp)
    if err != nil {
        glog.Flush()
        glog.Fatalf("Could not read: %v", err)
    }
    if *rndm {
        for i := range workers {
            j := rand.Intn(i + 1)
            glog.V(2).Infof("i %v; j %v", i, j)
            glog.V(2).Infof("swapping %v and %v", workers[i], workers[j])
            workers[i], workers[j] = workers[j], workers[i]
        }
    }
    glog.V(2).Infof("got workers %v", workers)
    buddies := makeBuddies(workers, grpSz)
    glog.V(2).Infof("got buddy groups %v", buddies)
    for i, b := range buddies {
        fmt.Printf("Group %v: %v\n", i, b)
    }
    glog.Flush()
}