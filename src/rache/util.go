package rache

import(
  "time"
  "strconv"
)

var DAYS = map[string]string{
  "Sunday":"128",
  "Monday":"64",
  "Tuesday":"32",
  "Wednesday":"16",
  "Thursday":"8",
  "Friday":"4",
  "Saturday":"2",
}

func UnixTimeToSegment(ts string) (day string, m int){
  t ,_:= strconv.ParseInt(ts,10,64)
  tdate := time.Unix(t,0)
  day = DAYS[tdate.Weekday().String()]
  m = HourToMinutes(tdate.Hour()) + tdate.Minute()
  return
}

func HourToMinutes(hour int) (m int){
  m = 60 * hour
  return
}

func TimeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  Logger.Infof("%s took %s",name,elapsed)
}

func RandomDestinationProperty() map[string]string{
  return map[string]string{
    "destination":"destination",
    "destination_title":"Jernz",
    "destinatation_attr":"",
    "call_type":"",
    "destination_property_name":"",
    "recording_percentage":"",
    "transfer_pattern":"",
    "transfer_method":"",
    "transfer_type":"",
    "outcome_timeout":"",
    "retry_count":"",
    "terminate":"",
    "max_speed_digits":"",
    "music_on_hold":"",
    "dial_or_block":"",
    "pass_parentcallID":"",
    "cdr_auth":"",
    "outdial_format":"",
    "agent_type":"",
    "delay_recording":"",
    "latched_recording":"",
    "hear_dtmf":"",
    "commands_ok":"",
    "dtmf_from_o":"",
    "dtmf_to_o":"",
    "dest_loc":"",
    "isup_enabled":"",
    "target_ack":"",
    "pass_ids":"",
    "pass_orig_dnis":"",
    "pass_ced":"",
    "delay_call_establish":"",
    "substitute_variables":"",
    "route_on_refer_user_id":"",
    "process_accumulated_xfer_digits":"",
    "initial_page_uri_xheader":"",
    "block_session_progress_inbound":"",
    "notifications_enabled":"",
    "early_media":"",
    "destination_attribute_bits":"",
    "queue_cti":"",
    "cti_routine":"",
    "vail_command_mask0":"",
    "vail_command_mask1":"",
    "privacy":"",
    "ani_override":"",
    "network_type":"",
  }
}
