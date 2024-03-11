/* eslint-disable no-unused-vars */
/*

CRUD
poolBuyer: isys2099_group9_app_buyer_user

query for creating ib_ord: put
CALL sp_place_buyer_order(?, ?, ?, (OUT?)), []

query retrive: get
- SELECT * FROM buyer_order;

query update:
- UPDATE buyer_order SET ??

query delete: delete
- DELETE buyer_order WHERE ??

For buyer orders, Buyer can: 

query accpet: post
UPDATE buyer_order SET order_status = 'A'

query reject: post
UPDATE buyer_order SET order_status = 'R'

Errors handling:

    -1 on rollback
    0 on successful commit
    1 on not enough warehouse space to return product
    2 on buyer_order_id not exist

err_msg = concat('Order already accepted.');
err_msg = concat('Order already rejected.');
err_msg = concat('Order already pending.');
err_msg = 'Cannot reopen rejected buyer order. Place a new order with sp_return_product_from_buyer_order() instead.';
err_msg = concat('Cannot accept rejected buyer order. Place a new order with sp_return_product_from_buyer_order() instead.');
err_msg = 'Cannot reopen accepted buyer order. Place a new order with sp_return_product_from_buyer_order() instead.';
err_msg = concat('Cannot reject accepted buyer order.');

*/

const { db, model } = require("../models");

// Buyer Order

/*
 sp_place_buyer_order(order_quantity: int, order_product_id: int, buyer_username: varchar(45), OUT result: int)

 OUT result:
    -1 on rollback
    0 on successful commit
    1 on not enough stockpile
    2 on product_id or buyer_id not exist
 */
const placeOrder = async (req, res) => {
  try {
    const order_product_id = req.params.id;
    const buyer_username = req.username;
    let { order_quantity } = req.body;

    // Double check if the product is still instock
    const [productInstock] = await db.poolBuyer.query(
      "SELECT quantity FROM stockpile WHERE product_id = ?",
      [order_product_id]
    );

    if (!productInstock[0]) {
      return res.status(200).json({ message: "Product is out of stock!" });
    }

    if (productInstock[0].quantity < order_quantity) {
      return res.status(200).json({ message: "Not enough product instock" });
    }

    // Check if the product id is already in the buyer order database
    const [existingOrders] = await db.poolBuyer.query(
      "SELECT * FROM buyer_order WHERE product_id = ? AND buyer = ?",
      [order_product_id, buyer_username],
    );

    if (existingOrders.length > 0) {
      // Get the product quantity of the original buyer order and add to the new buyer order quantity
      const existingOrder = existingOrders[0];
      order_quantity += existingOrder.quantity;

      // Call sp_return_product_from_buyer_order to return the product quantity of the original order
      await db.poolBuyer.query("CALL sp_return_product_from_buyer_order(?, @result)", [existingOrder.id]);
      const [[{ result: returnResultCode }]] = await db.poolBuyer.query("SELECT @result as result");
      if (returnResultCode !== 0) {
        return res.status(500).json({ error: "An error occurred while returning the product quantity", result: returnResultCode });
      }

      // Delete the original buyer order
      await db.poolBuyer.query("DELETE FROM buyer_order WHERE id = ?", [existingOrder.id]);
    }

    // Call the stored procedure to place an order
    await db.poolBuyer.query("CALL sp_place_buyer_order(?, ?, ?, @result)", [
      order_quantity,
      order_product_id,
      buyer_username,
    ]);
    const [[{ result: resultCode }]] = await db.poolBuyer.query(
      "SELECT @result as result",
    );
    if (resultCode === 0) {
      return res
        .status(200)
        .json({ message: "Order placed successfully", result: resultCode });
    } else if (resultCode === 1) {
      return res
        .status(400)
        .json({ error: "Not enough stockpile", result: resultCode });
    } else if (resultCode === 2) {
      return res
        .status(400)
        .json({ error: "Product or buyer does not exist", result: resultCode });
    }
    return res.status(500).json({
      error: "An error occurred while processing your request",
      result: resultCode,
    });
  } catch (error) {
    res.status(500).json({ error: "This is an error" });
  }
};



const getAllBuyerOrders = async (req, res) => {
  try {
    const [results] = await db.poolBuyer.query(`
        SELECT buyer_order.*, product.title AS product_title, product.category AS category, product.price AS price
        FROM buyer_order JOIN product ON buyer_order.product_id = product.id
        `);
    return res.json(results);
  } catch (error) {
    console.error("error: " + error.stack);
    return res.status(500).json({ error: "Internal server error" });
  }
};

const getBuyerOrderByID = async (req, res) => {
  try {
    let buyerOrderID = req.params.id;
    const [results] = await db.poolBuyer.query(
      `
            SELECT buyer_order.*, product.title AS product_title, product.category AS category, product.price AS price
            FROM buyer_order JOIN product ON buyer_order.product_id = product.id
            WHERE buyer_order.id = ?
        `,
      [buyerOrderID],
    );
    if (results.length === 0) {
      return res
        .status(404)
        .json({ error: `Buyer order with id: ${buyerOrderID} not found` });
    }
    return res.json(results);
  } catch (error) {
    console.error("error: " + error.stack);
    return res.status(500).json({ error: "Internal server error" });
  }
};

const getBuyerOrderByCategory = async (req, res) => {
  try {
    let category = req.params.category;
    const [results] = await db.poolBuyer.query(
      `
            SELECT buyer_order.*, product.title AS product_title, product.category AS category, product.price AS price
            FROM buyer_order JOIN product ON buyer_order.product_id = product.id
            WHERE product.category = ?
        `,
      [category],
    );
    return res.json(results);
  } catch (error) {
    console.error("error: " + error.stack);
    return res.status(500).json({ error: "Internal server error" });
  }
};

const getBuyerOrderByStatus = async (req, res) => {
  try {
    let status = req.params.status;
    const [results] = await db.poolBuyer.query(
      `
            SELECT buyer_order.*, product.title AS product_title, product.category AS category, product.price AS price
            FROM buyer_order JOIN product ON buyer_order.product_id = product.id
            WHERE order_status = ?
        `,
      [status],
    );
    return res.json(results);
  } catch (error) {
    console.error("error: " + error.stack);
    return res.status(500).json({ error: "Internal server error" });
  }
};

// PUT request for editting quantity of orders
const updateBuyerOrderQuantity = async (req, res) => {
  try {
    let buyerOrderID = req.params.id;
    let { quantity } = req.body;

    // Fetch the existing order
    const [orderResults] = await db.poolBuyer.query(
      `
            SELECT * FROM buyer_order WHERE id = ?
        `,
      [buyerOrderID],
    );
    if (orderResults.length === 0) {
      return res
        .status(404)
        .json({ error: `Buyer order with id: ${buyerOrderID} not found` });
    }

    const order = orderResults[0];
    const buyer_username = order.buyer;
    const order_product_id = order.product_id;

    // Call the stored procedure to place a new order
    await db.poolBuyer.query("CALL sp_place_buyer_order(?, ?, ?, @result)", [
      quantity,
      order_product_id,
      buyer_username,
    ]);
    const [[{ result: resultCode }]] = await db.poolBuyer.query(
      "SELECT @result as result",
    );
    if (resultCode !== 0) {
      return res.status(500).json({
        error: "An error occurred while placing the new order",
        result: resultCode,
      });
    }

    // Call the stored procedure to return the quantity of the original order
    await db.poolBuyer.query("CALL sp_return_product_from_buyer_order(?, @result)", [
      buyerOrderID,
    ]);
    const [[{ result: returnResultCode }]] = await db.poolBuyer.query(
      "SELECT @result as result",
    );
    if (returnResultCode !== 0) {
      return res.status(500).json({
        error: "An error occurred while returning the quantity of the original order",
        result: returnResultCode,
      });
    }

    // Delete the original order
    const [deleteResults] = await db.poolBuyer.query(
      `
              DELETE FROM buyer_order WHERE id = ?
          `,
      [buyerOrderID],
    );

    return res
      .status(200)
      .json({ message: "Order updated successfully", result: resultCode });
  } catch (error) {
    console.error("error: " + error.stack);
    return res.status(500).json({ error: "Internal server error" });
  }
};

const updateBuyerOrderStatusAccept = async (req, res) => {
  try {
    let buyerOrderID = req.params.id;
    const [results] = await db.poolBuyer.query(
      `
            UPDATE buyer_order SET order_status = 'A', fulfilled_date = DATE(SYSDATE()), fulfilled_time = TIME(SYSDATE()) WHERE id = ?
        `,
      [buyerOrderID],
    );
    if (results.affectedRows === 0) {
      return res
        .status(404)
        .json({ error: `Buyer order with id: ${buyerOrderID} not found` });
    }
    return res.json({
      message: `Status of buyer order with id: ${buyerOrderID} updated to accepted`,
    });
  } catch (error) {
    console.error("error: " + error.stack);
    return res.status(500).json({ error: "Internal server error" });
  }
};

const updateBuyerOrderStatusReject = async (req, res) => {
  try {
    let buyerOrderID = req.params.id;
    const [results] = await db.poolBuyer.query(
      `
            UPDATE buyer_order SET order_status = 'R' WHERE id = ?
        `,
      [buyerOrderID],
    );
    if (results.affectedRows === 0) {
      return res
        .status(404)
        .json({ error: `Buyer order with id: ${buyerOrderID} not found` });
    }
    return res.json({
      message: `Status of buyer order with id: ${buyerOrderID} updated to rejected`,
    });
  } catch (error) {
    console.error("error: " + error.stack);
    return res.status(500).json({ error: "Internal server error" });
  }
};

const deleteBuyerOrder = async (req, res) => {
  const { id } = req.params;
  try {
    // Get the order_status of the order based on its id
    const [orderResults] = await db.poolBuyer.query("SELECT order_status FROM buyer_order WHERE id = ?", [id]);
    if (orderResults.length === 0) {
      return res.status(404).json({ error: `Buyer order with id: ${id} not found` });
    }

    const order_status = orderResults[0].order_status;
    console.log(order_status);

    // If the order status is pending ('P'), call sp_return_product_from_buyer_order to return that product quantity back
    if (order_status === 'P') {
      await db.poolBuyer.query("CALL sp_return_product_from_buyer_order(?, @result)", [id]);
      const [[{ result: resultCode }]] = await db.poolBuyer.query("SELECT @result as result");
      if (resultCode !== 0) {
        return res.status(500).json({ error: "An error occurred while returning the product quantity" });
      }
    }

    // Delete that buyer order based on its id
    await db.poolBuyer.query("DELETE FROM buyer_order WHERE id = ?", [id]);
    res.status(200).json({ message: `Buyer order with ID: ${id} deleted`, id: id });
  } catch (error) {
    console.error(error);
    res.status(500).json({ error: "An error occurred while deleting a buyer order" });
  }
};


module.exports = {
  placeOrder,
  getAllBuyerOrders,
  getBuyerOrderByID,
  getBuyerOrderByCategory,
  getBuyerOrderByStatus,
  updateBuyerOrderQuantity,
  updateBuyerOrderStatusAccept,
  updateBuyerOrderStatusReject,
  deleteBuyerOrder,
};
