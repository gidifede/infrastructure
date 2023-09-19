from datetime import datetime, timezone



def get_timestamp():
    
    return  datetime.now(timezone.utc).isoformat().replace('+00:00', 'Z')
    