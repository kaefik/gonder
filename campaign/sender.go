package campaign

import (
	"github.com/supme/gonder/models"
	"sync"
	"time"
)

type campaign struct {
	id, from_email, from_name, subject, body string
	sendUnsubscribe                          bool
	profileId, resendDelay, resendCount      int
	attachments                              []Attachment
}

func (c campaign) run(id string) {
	c.get(id)
	c.send()
	c.resend()
}

// Send campaign
func (c campaign) send() {
	var r recipient
	var wg sync.WaitGroup

	query, err := models.Db.Prepare("SELECT `id`, `email`, `name` FROM recipient WHERE campaign_id=? AND removed=0 AND status IS NULL")
	checkErr(err)
	defer query.Close()

	q, err := query.Query(c.id)
	checkErr(err)
	defer q.Close()

	for q.Next() {
		err = q.Scan(&r.id, &r.to_email, &r.to_name)
		checkErr(err)
		if r.unsubscribe(c.id) == false || c.sendUnsubscribe {
			models.Db.Exec("UPDATE recipient SET status='Sending', date=NOW() WHERE id=?", r.id)
			pid, iface, host := models.ProfileNext(c.profileId)
			go func(d recipient, p int, i, h string) {
				wg.Add(1)
				rs := d.send(&c, i, h)
				wg.Done()
				models.ProfileFree(p)
				models.Db.Exec("UPDATE recipient SET status=?, date=NOW() WHERE id=?", rs, d.id)
			}(r, pid, iface, host)
		} else {
			models.Db.Exec("UPDATE recipient SET status='Unsubscribe', date=NOW() WHERE id=?", r.id)
			camplog.Printf("Recipient id %s email %s is unsubscribed", r.id, r.to_email)
		}
	}
	wg.Wait()
	camplog.Printf("Done campaign id %s.", c.id)
}

func (c campaign) resend() {
	var r recipient

	count := c.countSoftBounce()
	if count == 0 {
		return
	}

	if c.resendCount != 0 {
		camplog.Printf("Start %d resend by campaign id %s ", count, c.id)
	}

	for n := 0; n < c.resendCount; n++ {
		if c.countSoftBounce() == 0 {
			return
		}
		time.Sleep(time.Duration(c.resendDelay) * time.Second)
		//query, err := models.Db.Prepare("SELECT DISTINCT r.`id`, r.`email`, r.`name` FROM `recipient` as r,`status` as s WHERE r.`campaign_id`=? AND r.`removed`=0 AND s.`bounce_id`=2 AND UPPER(`r`.`status`) LIKE CONCAT(\"%\",s.`pattern`,\"%\")")
		query, err := models.Db.Prepare("SELECT `id`, `email`, `name` FROM `recipient` WHERE `campaign_id`=? AND `removed`=0 AND LOWER(`status`) REGEXP '^((4[0-9]{2})|(dial tcp)|(read tcp)|(proxy)|(eof)).+'")
		checkErr(err)
		defer query.Close()

		q, err := query.Query(c.id)
		checkErr(err)
		defer q.Close()

		for q.Next() {
			err = q.Scan(&r.id, &r.to_email, &r.to_name)
			checkErr(err)
			if r.unsubscribe(c.id) == false || c.sendUnsubscribe {
				models.Db.Exec("UPDATE recipient SET status='Sending', date=NOW() WHERE id=?", r.id)
				p, i, h := models.ProfileNext(c.profileId)
				rs := r.send(&c, i, h)
				models.ProfileFree(p)
				models.Db.Exec("UPDATE recipient SET status=?, date=NOW() WHERE id=?", rs, r.id)
			} else {
				models.Db.Exec("UPDATE recipient SET status='Unsubscribe', date=NOW() WHERE id=?", r.id)
				camplog.Printf("Recipient id %s email %s is unsubscribed", r.id, r.to_email)
			}
		}
	}
}

func (c *campaign) countSoftBounce() int {
	var count int
	//err := models.Db.QueryRow("SELECT COUNT(DISTINCT r.`id`) FROM `recipient` as r,`status` as s WHERE r.`campaign_id`=? AND r.`removed`=0 AND s.`bounce_id`=2 AND UPPER(`r`.`status`) LIKE CONCAT('%',s.`pattern`,'%')", c.id).Scan(&count)
	err := models.Db.QueryRow("SELECT COUNT(`id`) FROM `recipient` WHERE `campaign_id`=? AND `removed`=0 AND LOWER(`status`) REGEXP '^((4[0-9]{2})|(dial tcp)|(proxy)|(eof)).+'", c.id).Scan(&count)
	if err != nil {
		camplog.Print(err)
	}
	return count
}

// Get all info for campaign
func (c *campaign) get(id string) {
	err := models.Db.QueryRow("SELECT t1.`id`,t3.`email`,t3.`name`,t1.`subject`,t1.`body`,t2.`id`,t1.`send_unsubscribe`,t2.`resend_delay`,t2.`resend_count` FROM `campaign` t1 INNER JOIN `profile` t2 ON t2.`id`=t1.`profile_id` INNER JOIN `sender` t3 ON t3.`id`=t1.`sender_id` WHERE t1.`id`=?", id).Scan(
		&c.id,
		&c.from_email,
		&c.from_name,
		&c.subject,
		&c.body,
		&c.profileId,
		&c.sendUnsubscribe,
		&c.resendDelay,
		&c.resendCount,
	)

	checkErr(err)

	attachment, err := models.Db.Prepare("SELECT `path`, `file` FROM attachment WHERE campaign_id=?")
	checkErr(err)
	defer attachment.Close()

	attach, err := attachment.Query(c.id)
	checkErr(err)
	defer attach.Close()

	c.attachments = nil
	var location string
	var name string
	for attach.Next() {
		err = attach.Scan(&location, &name)
		checkErr(err)
		c.attachments = append(c.attachments, Attachment{Location: location, Name: name})
	}
}
