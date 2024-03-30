import logging
from kiteconnect import KiteTicker


# Initialise
kws = KiteTicker("246y3a7zlg83xv2l", "k67M1LOdxtG8sNiL2hc9Murbk75L8co1")

def on_ticks(ws, ticks):
    # Callback to receive ticks.
    print(ticks)

def on_connect(ws, response):

    ws.subscribe([738561, 5633])

    # Set RELIANCE to tick in `full` mode.
    ws.set_mode(ws.MODE_LTP, [738561, 5633])

def on_close(ws, code, reason):
    # On connection close stop the main loop
    # Reconnection will not happen after executing `ws.stop()`
    ws.stop()

# Assign the callbacks.
kws.on_ticks = on_ticks
kws.on_connect = on_connect
kws.on_close = on_close

# Infinite loop on the main thread. Nothing after this will run.
# You have to use the pre-defined callbacks to manage subscriptions.
kws.connect()
