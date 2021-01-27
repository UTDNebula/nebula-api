async function _getAll(collection) {
    const all = await collection.get();
    const result = [];
    all.forEach(c => {
        result.push(c.data());
    });
    return result;
}

async function _post(collection, updated, increment, key_ref) {
    // get counter to set new id
    const counter = collection.doc('_counter');
    const count = await counter.get();
    const newId = count.data()['count'];
    updated['id'] = parseInt(newId);

    await collection.doc(updated[key_ref]).set(updated);
    await counter.update({ count: increment });
}

async function _findExact(collection, key, value) {
    console.log(key + ", " + value);
    const res = await collection.where(key, '==', value).get();
    if (res.empty) {
        return {};
    } else {
        console.log('done');
        return res.docs[0].data();
    }
}

async function _findFuzzy(collection, key, value) {
    const upper = value.replace(/.$/, c => String.fromCharCode(c.charCodeAt(0) + 1));
    const snapshot = await collection.where(key, '>=', value)
        .where(key, '<', upper)
        .get();
    if (snapshot.empty) {
        return [];
    } else {
        const result = [];
        snapshot.forEach(doc => {
            result.push(doc.data());
        });
        return result;
    }
}

async function _deleteById(collection, key, value, decrement) {
    const result = collection.where(key, '==', value);
    result.get().then(async (snapshot) => {
        if (snapshot.empty) {
            return false;
        } else {
            snapshot.docs[0].ref.delete();
            const counter = collection.doc('_counter');
            await counter.update({ count: decrement });
            return true;
        }
    });
}

async function _patch(collection, id, updated) {
    collection.where('id', '==', id).get()
        .then(snapshot => {
            if (snapshot.empty) {
                return false;
            } else {
                snapshot.docs[0].ref.update(updated);
                return true;
            }
        });
}

module.exports = {
    _getAll, _post, _findExact, _findFuzzy, _deleteById, _patch
}
