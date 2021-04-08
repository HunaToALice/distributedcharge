package accessor

type Client struct {
}

func (c *Client) ChargePart(t *Transaction, uuid string) {
}

func (c *Client) Commit(eventno string, uuid string, iscommit bool) {
}

func (c *Client) ReportResult(eventno string, result bool) {

}
