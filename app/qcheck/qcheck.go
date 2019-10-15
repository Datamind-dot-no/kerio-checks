package qcheck

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/datamind-dot-no/kerio-checks/config"

	"github.com/datamind-dot-no/kerio-checks/app/notifications"

	"github.com/karrick/godirwalk"
)

type Qchk struct {
	kerioChkConf  *config.Config
	notifications *notifications.Notifications
}

func New(c *config.Config, n *notifications.Notifications) *Qchk {
	return &Qchk{
		kerioChkConf:  c,
		notifications: n,
	}
}

func (q *Qchk) CheckQ() error {
	// use faster implementation for counting files as recommended at https://boyter.org/2018/03/quick-comparison-go-file-walk-implementations/
	count := 0
	fmt.Println(q.kerioChkConf.KerioStorePath + q.kerioChkConf.QueuePath)
	godirwalk.Walk(q.kerioChkConf.KerioStorePath+q.kerioChkConf.QueuePath, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			// we're counting the files with extension eml as those are the actual messages
			if strings.HasSuffix(osPathname, "eml") {
				count++

				// debug statement:
				//fmt.Printf("%s %s\n", de.ModeType(), osPathname)
			}

			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			// Your program may want to log the error somehow.
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)

			return godirwalk.SkipNode
		},
	})
	QueueLength := count
	fmt.Println("The number of messages in the queue is: " + strconv.Itoa(QueueLength))

	if QueueLength > q.kerioChkConf.QueueWarnThreshold {
		q.notifications.SendNotification(QueueLength)
	}
	return nil
}
