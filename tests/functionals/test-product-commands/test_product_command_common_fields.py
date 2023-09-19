from commands.product import *
from queries.product import ProductStatus, ProductDetails, ProductLocation, ProductTrack
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product, Empty, Carrier
import time
import utils
import pytest

def arrange_cmd_fields(product_id):
    return Product(product_id, 'Poste Delivery Business', 'BOX', []), \
        Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', []), \
        Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', []), \
        Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', []), \
        Carrier('motorbike', 'Postman12345', 'FE764OU')
    

@pytest.mark.functional
def test_product_command_common_fields_empty_product_KO(config,product_id):
    ## Arrange
    _, receiver, sender, location, _ = arrange_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=Empty(),
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
    assert cmds[0].payload['product'] == {}
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_common_fields_bad_product_KO(config,product_id):
    ## Arrange
    _, receiver, sender, location, _ = arrange_cmd_fields(product_id=product_id)
    product_unknown = Product(product_id, 'Poste Delivery Business', 'UNKNOWN', [])
    product_bad_attr_number = Product(product_id, 'Poste Delivery Business', 'PHARMA', 1)
    product_bad_attr_string = Product(product_id, 'Poste Delivery Business', 'BOX', 'bad string')
    cmds = [
        AcceptProduct(
                product=product_unknown,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender,
                location=location,
                attributes=[]
            ),
        AcceptProduct(
                product=product_bad_attr_number,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender,
                location=location,
                attributes=[]
            ),
        AcceptProduct(
                product=product_bad_attr_string,
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
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400

@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-503")
def test_product_command_common_fields_empty_receiver_KO(config,product_id):
    ## Arrange
    product, _, sender, location, _ = arrange_cmd_fields(product_id=product_id)
    cmds = [ 
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=Empty(),
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
    assert cmds[0].payload['receiver'] == {}
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_common_fields_bad_receiver_KO(config,product_id):
    ## Arrange
    product, _, sender, location, _ = arrange_cmd_fields(product_id=product_id)
    receiver_bad_zipcode = Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', 'a code', '3333333333', 'mario@rossi.com', 'Carrier\'s note', [])
    receiver_bad_attr_number = Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', 654)
    receiver_bad_attr_string = Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', 'bad attributes')
    cmds = [
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver_bad_zipcode,
                sender=sender,
                location=location,
                attributes=[]
            ),
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver_bad_attr_number,
                sender=sender,
                location=location,
                attributes=[]
            ),
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver_bad_attr_string,
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
def test_product_command_common_fields_empty_sender_KO(config,product_id):
    ## Arrange
    product, receiver, _, location, _ = arrange_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=Empty(),
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
    assert cmds[0].payload['sender'] == {}
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_common_fields_bad_sender_KO(config,product_id):
    ## Arrange
    product, receiver, _, location, _ = arrange_cmd_fields(product_id=product_id)
    sender_bad_zipcode = Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', 'bad code', '2222222222', [])
    sender_bad_attr_number = Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', 45)
    sender_bad_attr_string = Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', 'bad string')
    cmds = [
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender_bad_zipcode,
                location=location,
                attributes=[]
            ),
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender_bad_attr_number,
                location=location,
                attributes=[]
            ),
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender_bad_attr_string,
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
def test_product_command_common_fields_empty_location_KO(config,product_id):
    ## Arrange
    product, receiver, sender, _, _ = arrange_cmd_fields(product_id=product_id)
    cmds = [
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender,
                location=Empty(),
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
    assert cmds[0].payload['location'] == {}
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_common_fields_bad_location_KO(config,product_id):
    ## Arrange
    product, receiver, sender, _, _ = arrange_cmd_fields(product_id=product_id)
    location_bad_zipcode = Location('UP', 'Roma', 'Viale Europa 190', 'bad_one', 'Italia', '55Y90', [])
    location_bad_attr_number = Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', 146)
    location_bad_attr_string = Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', 'a bad string')
    cmds = [
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender,
                location=location_bad_zipcode,
                attributes=[]
            ),
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender,
                location=location_bad_attr_number,
                attributes=[]
            ),
        AcceptProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                receiver=receiver,
                sender=sender,
                location=location_bad_attr_string,
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
def test_product_command_common_fields_empty_carrier_KO(config,product_id):
    ## Arrange
    product, _, _, _, _ = arrange_cmd_fields(product_id=product_id)
    cmds = [ 
        StartProductDelivery(
                product=product,
                carrier=Empty(),
                timestamp=utils.get_timestamp(),
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
    assert cmds[0].payload['carrier'] == {}
    assert len(responses) == 1
    assert responses[0].status_code == 400