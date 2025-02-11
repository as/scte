# scte
Package scte implements a SCTE_35 packet reader based on the following standard:

https://dutchguild.nl/event/13/attachments/82/203/SCTE_35_2023r1.pdf

It decodes all non-deprecated commands and produces an informative data structure on the
command and any optional segment descriptors attached to it. SCTE-35 packets are used
primarily to signal AD-breaks and content decisioning. They are carried in MPEG transport
streams as well as DASH and HLS media playlists in base64 or hexidecimal formats.

# Installation

This repository comes with a Go library to parse SCTE-35 messages as well as a
command line executable to parse them from standard input. To install both:

```
go get github.com/as/scte
go install github.com/as/scte/cmd/scte@latest
```

# Library

```
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/as/scte"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	p, err := scte.Parse(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(p)
}
```

# Command Line Examples

# 14.1. time_signal Placement Opportunity Start
```
echo FC3034000000000000FFFFF00506FE72BD0050001E021C435545494800008E7FCF0001A599B00808000000002CA0A18A3402009AC9D17E | scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":52,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":255,"Tier":4095,"CmdLen":5,"CmdType":6,"CmdName":"time_signal","Cmd":{"HasPTS":true,"Res":63,"PTS":1924989008},"DescLen":30,"Desc":[{"Tag":2,"Len":28,"ID":1129661769,"TagName":"segmentation_descriptor","Error":"warning: sub_segment_num and sub_segments_expected are required but missing","EventID":1207959694,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":true,"DeliveryUnrestricted":false,"WebDelivery":false,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":27630000,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACygoYo=","SegType":52,"SegTypeName":"ProviderPlacementOpportunityStart","SegNum":2,"SegExp":0,"SubSegNum":0,"SubSegExp":0}],"Stuffing":null,"ECRC32":0,"CRC32":2596917630}
```

# 14.2. splice_insert
```
echo '/DAvAAAAAAAA///wFAVIAACPf+/+c2nALv4AUsz1AAAAAAAKAAhDVUVJAAABNWLbowo='| scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":47,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":255,"Tier":4095,"CmdLen":20,"CmdType":5,"CmdName":"splice_insert","Cmd":{"Cmd":null,"EventID":1207959695,"Cancel":false,"Res0":127,"OutOfNetwork":true,"HasSplice":true,"HasDuration":true,"Immediate":false,"Compliant":true,"Res1":7,"Time":{"HasPTS":true,"Res":63,"PTS":1936310318},"AutoReturn":true,"Res2":63,"BreakDur":5426421,"ProgID":0,"Avail":0,"AvailExp":0},"DescLen":10,"Desc":[{"Tag":0,"Len":8,"ID":1129661769,"Data":"AAABNQ=="}],"Stuffing":null,"ECRC32":0,"CRC32":1658561290}

```
# 14.3. time_signal Placement Opportunity End
```
echo 'FC302F000000000000FFFFF00506FE746290A000190217435545494800008E7F9F0808000000002CA0A18A350200A9CC6758'| scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":47,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":255,"Tier":4095,"CmdLen":5,"CmdType":6,"CmdName":"time_signal","Cmd":{"HasPTS":true,"Res":63,"PTS":1952616608},"DescLen":25,"Desc":[{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959694,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACygoYo=","SegType":53,"SegTypeName":"ProviderPlacementOpportunityEnd","SegNum":2,"SegExp":0,"SubSegNum":0,"SubSegExp":0}],"Stuffing":null,"ECRC32":0,"CRC32":2848745304}

```
# 14.4. time_signal Program Start/End
```
echo '/DBIAAAAAAAA///wBQb+ek2ItgAyAhdDVUVJSAAAGH+fCAgAAAAALMvDRBEAAAIXQ1VFSUgAABl/nwgIAAAAACyk26AQAACZcuND'| scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":72,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":255,"Tier":4095,"CmdLen":5,"CmdType":6,"CmdName":"time_signal","Cmd":{"HasPTS":true,"Res":63,"PTS":2051901622},"DescLen":50,"Desc":[{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959576,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACzLw0Q=","SegType":17,"SegTypeName":"ProgramEnd","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0},{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959577,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACyk26A=","SegType":16,"SegTypeName":"ProgramStart","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0}],"Stuffing":null,"ECRC32":0,"CRC32":2574443331}

```
# 14.5. time_signal Program Overlap Start
```
echo 'FC302F000000000000FFFFF00506FEAEBFFF640019021743554549480000087F9F0808000000002CA56CF5170000951DB0A8'| scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":47,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":255,"Tier":4095,"CmdLen":5,"CmdType":6,"CmdName":"time_signal","Cmd":{"HasPTS":true,"Res":63,"PTS":2931818340},"DescLen":25,"Desc":[{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959560,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACylbPU=","SegType":23,"SegTypeName":"ProgramOverlapStart","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0}],"Stuffing":null,"ECRC32":0,"CRC32":2501750952}

```
# 14.6. time_signal Program Blackout Override / Program End
```
echo 'FC3048000000000000FFFFF00506FE932E380B00320217435545494800000A7F9F0808000000002CA0A1E3180000021743554549480000097F9F0808000000002CA0A18A110000B4217EB0'| scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":72,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":255,"Tier":4095,"CmdLen":5,"CmdType":6,"CmdName":"time_signal","Cmd":{"HasPTS":true,"Res":63,"PTS":2469279755},"DescLen":50,"Desc":[{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959562,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACygoeM=","SegType":24,"SegTypeName":"ProgramBlackoutOverride","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0},{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959561,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACygoYo=","SegType":17,"SegTypeName":"ProgramEnd","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0}],"Stuffing":null,"ECRC32":0,"CRC32":3022094000}

```
# 14.7. time_signal Program End
```
echo 'FC302F000000000000FFFFF00506FEAEF17C4C0019021743554549480000077F9F0808000000002CA56C97110000C4876A2E'| scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":47,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":255,"Tier":4095,"CmdLen":5,"CmdType":6,"CmdName":"time_signal","Cmd":{"HasPTS":true,"Res":63,"PTS":2935061580},"DescLen":25,"Desc":[{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959559,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACylbJc=","SegType":17,"SegTypeName":"ProgramEnd","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0}],"Stuffing":null,"ECRC32":0,"CRC32":3297208878}

```
# 14.8. time_signal Program Start/End - Placement Opportunity End
```
echo '/DBhAAAAAAAA///wBQb+qM1E7QBLAhdDVUVJSAAArX+fCAgAAAAALLLXnTUCAAIXQ1VFSUgAACZ/nwgIAAAAACyy150RAAACF0NVRUlIAAAnf58ICAAAAAAsstezEAAAihiGnw=='| scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":97,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":255,"Tier":4095,"CmdLen":5,"CmdType":6,"CmdName":"time_signal","Cmd":{"HasPTS":true,"Res":63,"PTS":2832024813},"DescLen":75,"Desc":[{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959725,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACyy150=","SegType":53,"SegTypeName":"ProviderPlacementOpportunityEnd","SegNum":2,"SegExp":0,"SubSegNum":0,"SubSegExp":0},{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959590,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACyy150=","SegType":17,"SegTypeName":"ProgramEnd","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0},{"Tag":2,"Len":23,"ID":1129661769,"TagName":"segmentation_descriptor","EventID":1207959591,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":false,"DeliveryUnrestricted":false,"WebDelivery":true,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":0,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACyy17M=","SegType":16,"SegTypeName":"ProgramStart","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0}],"Stuffing":null,"ECRC32":0,"CRC32":2316863135}
```

# Elemental OATCLS-SCTE35 - Placement Opportunity Start
```
echo '/DA0AAAAAAAAAAAABQb+ADAQ6QAeAhxDVUVJQAAAO3/PAAEUrEoICAAAAAAg+2UBNAAANvrtoQ== '| scte
{"Table":252,"SSI":false,"Priv":false,"SAP":3,"Len":52,"Ver":0,"Enc":false,"EncAlg":0,"PTSA":0,"CWI":0,"Tier":0,"CmdLen":5,"CmdType":6,"CmdName":"time_signal","Cmd":{"HasPTS":true,"Res":63,"PTS":3150057},"DescLen":30,"Desc":[{"Tag":2,"Len":28,"ID":1129661769,"TagName":"segmentation_descriptor","Error":"warning: sub_segment_num and sub_segments_expected are required but missing","EventID":1073741883,"Cancel":false,"Compliant":true,"Res0":63,"Segmented":true,"HasDuration":true,"DeliveryUnrestricted":false,"WebDelivery":false,"NoBlackout":true,"CanArchive":true,"Restrictions":3,"Res1":0,"Duration":18132042,"UPIDType":8,"UPIDLen":8,"UPID":"AAAAACD7ZQE=","SegType":52,"SegTypeName":"ProviderPlacementOpportunityStart","SegNum":0,"SegExp":0,"SubSegNum":0,"SubSegExp":0}],"Stuffing":null,"ECRC32":0,"CRC32":922414497}

```

# TODO

- Writing the Packets back into a bitstream
