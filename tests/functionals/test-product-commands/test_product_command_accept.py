from commands.product import *
from queries.product import ProductStatus, ProductDetails, ProductLocation, ProductTrack
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product, Empty
import time
import utils
import pytest

def arrange_accept_cmd_fields(product_id):
    return Product(product_id, 'Poste Delivery Business', 'BOX', []), \
        Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', []), \
        Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', []), \
        Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])


@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-501, LOI-502, LOI-533")
def test_product_command_accept_OK(config,product_id):
    ## Arrange
    product, receiver, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    ap = AcceptProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=receiver,
            sender=sender,
            location=location,
            attributes=[]
        )
    cmds = []
    cmds.append(ap) 
    
    queries = [
        ProductStatus(
            id = product_id
        ), 
        ProductDetails(
            id = product_id
        ), 
        ProductLocation(
            id = product_id
        ), 
        ProductTrack(
            id = product_id)
    ]
    
    expStatus = 'Accepted'

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()
    assert responses[0].status_code == 202
    time.sleep(5)
    results =  plan.start_query()

    ## Assert
    assert len(results) == 4
    # status check
    assert results[0]['product']['id'] == product_id
    assert results[0]['status'] == expStatus

    # details check
    assert results[1]['product'] == cmds[0].payload['product']
    assert results[1]['sender'] == cmds[0].payload['sender']
    assert results[1]['receiver'] == cmds[0].payload['receiver']

    # location check
    assert results[2]['product'] == cmds[0].payload['product']
    assert results[2]['location'] == cmds[0].payload['location']

    # track check
    assert results[3]['product']['id'] == product_id
    assert len(results[3]['trackList']) == 1
    assert results[3]['trackList'][0]['description'] == expStatus

@pytest.mark.functional
def test_product_command_accept_missing_product_KO(config,product_id):
    ## Arrange
    _, receiver, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender,
                location=location,
                attributes=[]
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['product']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_accept_missing_receiver_KO(config,product_id):
    ## Arrange
    product, _, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    cmds = [ 
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                sender=sender,
                location=location,
                attributes=[]
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['receiver']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_accept_missing_sender_KO(config,product_id):
    ## Arrange
    product, receiver, _, location = arrange_accept_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                location=location,
                attributes=[]
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['sender']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_accept_missing_location_KO(config,product_id):
    ## Arrange
    product, receiver, sender, _ = arrange_accept_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender,
                attributes=[]
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['location']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_accept_missing_timestamp_KO(config,product_id):
    ## Arrange
    product, receiver, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=product,
                receiver=receiver,
                sender=sender,
                location=location,
                attributes=[]
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['timestamp']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_accept_bad_timestamp_KO(config,product_id):
    ## Arrange
    product, receiver, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=product,
                timestamp=123456789,
                receiver=receiver,
                sender=sender,
                location=location,
                attributes=[]
            ),
        AcceptProduct(
                product=product,
                timestamp='A bad string',
                receiver=receiver,
                sender=sender,
                location=location,
                attributes=[]
            ),
        AcceptProduct(
                product=product,
                timestamp='2022-12-07T17:48:33',
                receiver=receiver,
                sender=sender,
                location=location,
                attributes=[]
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400

@pytest.mark.functional
def test_product_command_accept_missing_attributes_KO(config,product_id):
    ## Arrange
    product, receiver, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=product,
                receiver=receiver,
                sender=sender,
                location=location,
                timestamp=utils.get_timestamp(),
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['attributes']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_accept_bad_attributes_KO(config,product_id):
    ## Arrange
    product, receiver, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=product,
                receiver=receiver,
                sender=sender,
                location=location,
                timestamp=utils.get_timestamp(),
                attributes='a string'
            ),
        AcceptProduct(
                product=product,
                receiver=receiver,
                sender=sender,
                location=location,
                timestamp=utils.get_timestamp(),
                attributes=5
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400

@pytest.mark.functional
def test_product_command_accept_missing_api_key_KO(config,product_id):
    ## Arrange
    product, receiver, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    ap = AcceptProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=receiver,
            sender=sender,
            location=location,
            attributes=[]
        )
    cmds = []
    cmds.append(ap) 
    
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        # missign api key
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403
    
@pytest.mark.functional
def test_product_command_accept_wrong_api_key_KO(config,product_id):
    ## Arrange
    product, receiver, sender, location = arrange_accept_cmd_fields(product_id=product_id)
    ap = AcceptProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=receiver,
            sender=sender,
            location=location,
            attributes=[]
        )
    cmds = []
    cmds.append(ap) 
    
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key='fakeAPiKey'
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403
