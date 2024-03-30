import datetime
from enum import Enum
from datetime import timedelta        

class enum_historical_max(Enum):
    minute = "4mo"
    hour = "1y"
    week = "5y"
    oneminute = "2mo"
    

class get_max_history(enum_historical_max):
    
    def get_timedelta_from_enum(enum_value):
        timedelta_mapping = {
            "4mo": timedelta(weeks=20), # Approximately 4 months
            "1y": timedelta(weeks=52),   # Approximately 1 year
            "5y": timedelta(weeks=260), # Approximately 5 years
            "2mo": timedelta(weeks=8),   # Approximately 2 months
        }
        return timedelta_mapping.get(enum_value, timedelta())
    
    def max_history_date(get_timedelta_from_enum, time_period_enum):
        
        current_date = str(datetime.now())
        
        time_period_timedelta = get_timedelta_from_enum(time_period_enum.value)
        
        last_date = current_date - time_period_timedelta