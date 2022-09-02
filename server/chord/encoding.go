package chord

type TagEncoding struct {
	FileName      string `json:"file_name"`
	FileExtension string `json:"file_extension"`
	NodeID        []byte `json:"node_id"`
	NodeIP        string `json:"node_ip"`
	NodePort      string `json:"node_port"`
}
