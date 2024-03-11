import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import PropTypes from "prop-types";
import { getDataAPI, putDataAPI } from "../../../api/apiRequest";

const BuyerOrder = ({ compData, compFunction }) => {
  const { edited, editedBuyerOrderData } = compData;
  const { setPopUpState, setBuyerOrdersData } = compFunction;
  const [stockpile, setStockpile] = useState([]);
  const [buyerOrders, setBuyerOrders] = useState([]);

  const BuyerOrderState = {
    order_id: "",
    quantity: 0,
  };

  const [buyerOrderData, setBuyerOrderData] = useState(BuyerOrderState);
  const { quantity, order_id } = buyerOrderData;

  useEffect(() => {
    if (edited) {
      setBuyerOrderData(preData => ({
        ...preData,
        quantity: editedBuyerOrderData.quantity
          ? editedBuyerOrderData.quantity
          : "",
        order_id: editedBuyerOrderData.id ? editedBuyerOrderData.id : "",
      }));
    }

    const fetchDetailProduct = async () => {
      try {
        const response = await getDataAPI("stock");

        setStockpile(response.data);

      } catch (error) {
        console.error("Error fetching products:", error);
      }
    };

    fetchDetailProduct();
    console.log(stockpile);

    const fetchBuyerOrderByID = async () => {
      try {
        const response = await getDataAPI("buyer-order");

        setBuyerOrders(response.data);

      } catch (error) {
        console.error("Error fetching products:", error);
      }
    };

    fetchBuyerOrderByID();
    console.log(buyerOrders);

  }, [
    edited,
    editedBuyerOrderData.id,
    editedBuyerOrderData.order_id,
    editedBuyerOrderData.quantity,
  ]);

  const handleChangeInput = e => {
    const { name, value } = e.target;
    const parseValue = name === "quantity" ? parseInt(value, 10) : value;
    setBuyerOrderData(preState => ({ ...preState, [name]: parseValue }));
  };



  const handleSubmitData = async e => {
    e.preventDefault();

    try {
      // Find the product id in the buyer orders
      const productInOrder = buyerOrders.find(item => item.id === buyerOrderData.order_id);
      console.log(productInOrder);
      // Find the product in the stockpile
      const productInStockpile = stockpile.find(item => item.product_id === productInOrder.product_id);
      console.log(productInStockpile);
      // If the product doesn't exist in the stockpile or the requested quantity is more than the available quantity
      if (!productInStockpile || buyerOrderData.quantity > productInStockpile.quantity) {
        toast.error("Not enough quantity in stockpile");
        return;
      }

      if (buyerOrderData.quantity > 0) {
        if (edited) {
          const response = await putDataAPI(
            `buyer-order/${editedBuyerOrderData.id}/quantity`,
            { quantity: buyerOrderData.quantity },
          );
          if (response.status === 200 || response.status === 201) {
            setBuyerOrdersData(preData =>
              preData.map(obj => {
                if (obj.id === editedBuyerOrderData.id) {
                  return {
                    ...obj,
                    quantity: buyerOrderData.quantity,
                  };
                }
                return obj;
              }),
            );
            toast.success(`Edit order ${editedBuyerOrderData.id} successfully`);
          }
        }

        handleClosePopUpForm();
      } else {
        toast.error("Quantity of a product must be larger than 1");
      }
    } catch (error) {
      toast.error(error.response?.data?.error);
    }
  };

  const handleClosePopUpForm = () => {
    if (edited) {
      setPopUpState(prevState => ({
        ...prevState,
        state: !prevState.state,
        edited: false,
      }));
    }
  };

  return (
    <div className="popup_container p-4" style={{ top: "20%" }}>
      <form onSubmit={handleSubmitData}>
        <div className="container_fluid">
          <div className="row">
            <div className="col-6">
              <div className="mb-3">
                <label className="form-label">Product ID</label>
                <input
                  type="number"
                  className="form-control"
                  name="order_id"
                  value={order_id}
                  disabled={true}
                />
              </div>
            </div>
            <div className="col-6">
              <div className="mb-3">
                <label className="form-label">Quantity</label>
                <input
                  type="number"
                  className="form-control"
                  name="quantity"
                  value={quantity}
                  onChange={handleChangeInput}
                />
              </div>
            </div>
          </div>
          <div className="submit_btn">
            <span className="btn" onClick={handleClosePopUpForm}>
              Cancel
            </span>
            <button className="btn btn-outline-primary ms-2" type="submit">
              Edit
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

BuyerOrder.propTypes = {
  compData: PropTypes.shape({
    edited: PropTypes.bool.isRequired,
    editedBuyerOrderData: PropTypes.object,
  }).isRequired,

  compFunction: PropTypes.shape({
    setPopUpState: PropTypes.func.isRequired,
    setBuyerOrdersData: PropTypes.func.isRequired,
  }),
};

export default BuyerOrder;
